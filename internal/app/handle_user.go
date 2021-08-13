package app

import (
	"errors"
	"github.com/nebisin/api_structure/internal/store"
	"github.com/nebisin/api_structure/pkg/request"
	"github.com/nebisin/api_structure/pkg/response"
	"net/http"
	"strings"
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

	repo := store.NewUserRepository(s.db)

	if err := repo.Insert(user); err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			errs := map[string]string{"email": "is already exist"}
			response.FailedValidationResponse(w, errs)
		default:
			response.ServerErrorResponse(w, r, s.logger, err)
		}
		return
	}
	err = response.JSONResponse(w, http.StatusCreated, response.Envelope{"user": user});
	if err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
	}
}
