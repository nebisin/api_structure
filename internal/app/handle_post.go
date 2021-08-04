package app

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/nebisin/api_structure/internal/store"
	"github.com/nebisin/api_structure/pkg/request"
	"github.com/nebisin/api_structure/pkg/response"
	"net/http"
	"strconv"
)

func (s *server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string   `json:"title" validate:"required"`
		Body  string   `json:"body" validate:"required"`
		Tags  []string `json:"tags,omitempty" validate:"unique"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		response.BadRequestResponse(w, err)
		return
	}

	if err := request.ValidateInput(&input); err != nil {
		response.FailedValidationResponse(w, err)
		return
	}

	post := store.Post{
		Title:     input.Title,
		Body:      input.Body,
		Tags:      input.Tags,
	}

	repo := store.NewPostRepository(s.DB)

	err := repo.Insert(&post)
	if err != nil {
		response.ServerErrorResponse(w)
		return
	}

	if err := response.JSONResponse(w, http.StatusCreated, post); err != nil {
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

	repo := store.NewPostRepository(s.DB)

	post, err := repo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.NotFoundResponse(w)
		default:
			response.ServerErrorResponse(w)
		}
		return
	}

	if err := response.JSONResponse(w, http.StatusOK, post); err != nil {
		s.Logger.Println(err)
		response.ServerErrorResponse(w)
	}
}

func (s *server) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.NotFoundResponse(w)
		return
	}

	repo := store.NewPostRepository(s.DB)

	post, err := repo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.NotFoundResponse(w)
		default:
			response.ServerErrorResponse(w)
		}
		return
	}

	var input struct{
		Title string   `json:"title" validate:"required"`
		Body  string   `json:"body" validate:"required"`
		Tags  []string `json:"tags,omitempty" validate:"unique"`
	}

	if err = request.ReadJSON(w, r, &input); err != nil {
		response.BadRequestResponse(w, err)
	}

	if err := request.ValidateInput(&input); err != nil {
		response.FailedValidationResponse(w, err)
		return
	}

	post.Title = input.Title
	post.Body = input.Body
	post.Tags = input.Tags

	if 	err := repo.Update(post); err != nil {
		response.ServerErrorResponse(w)
		return
	}

	if err := response.JSONResponse(w, http.StatusOK, post); err != nil {
		response.ServerErrorResponse(w)
	}
}

func (s *server) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.NotFoundResponse(w)
		return
	}

	repo := store.NewPostRepository(s.DB)

	if 	err := repo.Delete(id); err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.NotFoundResponse(w)
		default:
			response.ServerErrorResponse(w)
		}
		return
	}

	err = response.JSONResponse(w, http.StatusOK, map[string]string{"message": "post successfully deleted"})
	if err != nil {
		response.ServerErrorResponse(w)
	}
}