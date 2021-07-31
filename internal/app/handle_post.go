package app

import (
	"github.com/gorilla/mux"
	"github.com/nebisin/api_structure/internal/store"
	"github.com/nebisin/api_structure/pkg/response"
	"github.com/nebisin/api_structure/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

func (s *server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string   `json:"title" validate:"required"`
		Body  string   `json:"body" validate:"required"`
		Tags  []string `json:"tags,omitempty" validate:"unique"`
	}


	if err := utils.ReadJSON(w, r, &input); err != nil {
		response.BadRequestResponse(w, err)
		return
	}

	if err := utils.ValidateInput(&input); err != nil {
		response.FailedValidationResponse(w, err)
		return
	}

	if err := response.JSONResponse(w, http.StatusCreated, input); err != nil {
		s.Logger.Println(err)
		response.ServerErrorResponse(w)
	}
}

func (s *server) handleShowPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.NotFoundResponse(w)
		return
	}

	post := store.Post{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Developing Rest API with Golang",
		Body:      "Lorem ipsum dolor sit amet.",
		Tags:      []string{"golang", "rest", "api"},
		Version:   1,
	}

	if err := response.JSONResponse(w, http.StatusOK, post); err != nil {
		s.Logger.Println(err)
		response.ServerErrorResponse(w)
	}
}
