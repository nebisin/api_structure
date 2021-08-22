package app

import (
	"errors"
	"github.com/nebisin/api_structure/internal/store"
	"github.com/nebisin/api_structure/pkg/request"
	"github.com/nebisin/api_structure/pkg/response"
	"net/http"
	"strings"
	"time"
)

func (s *server) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name" validate:"required,max=500"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,max=72,min=8"`
	}

	err := request.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, err)
	}

	if err := request.ValidateInput(&input); err != nil {
		response.FailedValidationResponse(w, err)
		return
	}

	user := &store.User{
		Name:      input.Name,
		Email:     strings.ToLower(input.Email),
		Activated: false,
	}

	if err := user.Password.Set(input.Password); err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
		return
	}

	if err := s.models.Users.Insert(user); err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			errs := map[string]string{"email": "is already exist"}
			response.FailedValidationResponse(w, errs)
		default:
			response.ServerErrorResponse(w, r, s.logger, err)
		}
		return
	}

	err = s.models.Permissions.AddForUser(user.ID, "posts:read")
	if err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
		return
	}

	token, err := s.models.Tokens.New(user.ID, 3*24*time.Hour, store.ScopeActivation)
	if err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
		return
	}

	s.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		if err := s.mailer.Send(user.Email, "user_welcome.tmpl", data); err != nil {
			s.logger.WithFields(map[string]interface{}{
				"request_method": r.Method,
				"request_url":    r.URL.String(),
			}).WithError(err).Error("background email error")
		}
	})

	err = response.JSONResponse(w, http.StatusAccepted, response.Envelope{"user": user})
	if err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
	}
}

func (s *server) handleActivateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlainText string `json:"token" validator:"required,max=26"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		response.BadRequestResponse(w, err)
		return
	}

	if err := request.ValidateInput(&input); err != nil {
		response.FailedValidationResponse(w, err)
		return
	}

	user, err := s.models.Users.GetForToken(store.ScopeActivation, input.TokenPlainText)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.FailedValidationResponse(w, map[string]string{"token": "invalid or expired activation token"})
		default:
			response.ServerErrorResponse(w, r, s.logger, err)
		}
		return
	}

	user.Activated = true

	if err := s.models.Users.Update(user); err != nil {
		switch {
		case errors.Is(err, store.ErrEditConflict):
			response.EditConflictResponse(w)
		default:
			response.ServerErrorResponse(w, r, s.logger, err)
		}
		return
	}

	if err := s.models.Tokens.DeleteAllForUser(store.ScopeActivation, user.ID); err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
		return
	}

	if err := response.JSONResponse(w, http.StatusOK, response.Envelope{"user": user}); err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
	}
}
