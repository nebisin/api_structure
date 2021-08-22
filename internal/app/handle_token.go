package app

import (
	"errors"
	"github.com/nebisin/api_structure/internal/store"
	"github.com/nebisin/api_structure/pkg/request"
	"github.com/nebisin/api_structure/pkg/response"
	"net/http"
	"time"
)

func (s *server) handleCreateAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,max=72,min=8"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		response.BadRequestResponse(w, err)
		return
	}

	if err := request.ValidateInput(&input); err != nil {
		response.FailedValidationResponse(w, err)
		return
	}

	user, err := s.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, s.logger, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
		return
	}

	if !match {
		response.InvalidCredentialsResponse(w, r)
		return
	}

	token, err := s.models.Tokens.New(user.ID, 24*time.Hour, store.ScopeAuthentication)
	if err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
		return
	}

	err = response.JSONResponse(w, http.StatusCreated, response.Envelope{"authentication_token": token})
	if err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
	}
}
