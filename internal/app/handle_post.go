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
		Title: input.Title,
		Body:  input.Body,
		Tags:  input.Tags,
	}

	repo := store.NewPostRepository(s.DB)

	err := repo.Insert(&post)
	if err != nil {
		response.ServerErrorResponse(w, s.Logger, err)
		return
	}

	if err := response.JSONResponse(w, http.StatusCreated, response.Envelope{"post": post}); err != nil {
		response.ServerErrorResponse(w, s.Logger, err)
	}
}

func (s *server) handleShowPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.NotFoundResponse(w, r)
		return
	}

	repo := store.NewPostRepository(s.DB)

	post, err := repo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, s.Logger, err)
		}
		return
	}

	if err := response.JSONResponse(w, http.StatusOK, response.Envelope{"post": post}); err != nil {
		response.ServerErrorResponse(w, s.Logger, err)
	}
}

func (s *server) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.NotFoundResponse(w, r)
		return
	}

	repo := store.NewPostRepository(s.DB)

	post, err := repo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, s.Logger, err)
		}
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(post.Version), 32) != r.Header.Get("X-Expected-Version") {
			response.EditConflictResponse(w)
			return
		}
	}

	var input struct {
		Title *string  `json:"title"`
		Body  *string  `json:"body"`
		Tags  []string `json:"tags,omitempty" validate:"unique"`
	}

	if err = request.ReadJSON(w, r, &input); err != nil {
		response.BadRequestResponse(w, err)
	}

	if err := request.ValidateInput(&input); err != nil {
		response.FailedValidationResponse(w, err)
		return
	}
	if input.Title != nil {
		post.Title = *input.Title
	}

	if input.Body != nil {
		post.Body = *input.Body
	}

	if input.Tags != nil {
		post.Tags = input.Tags
	}

	if err := repo.Update(post); err != nil {
		switch {
		case errors.Is(err, store.ErrEditConflict):
			response.EditConflictResponse(w)
		default:
			response.ServerErrorResponse(w, s.Logger, err)
		}
		return
	}

	if err := response.JSONResponse(w, http.StatusOK, response.Envelope{"post": post}); err != nil {
		response.ServerErrorResponse(w, s.Logger, err)
	}
}

func (s *server) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.NotFoundResponse(w, r)
		return
	}

	repo := store.NewPostRepository(s.DB)

	if err := repo.Delete(id); err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, s.Logger, err)
		}
		return
	}

	err = response.JSONResponse(w, http.StatusOK, response.Envelope{"message": "post successfully deleted"})
	if err != nil {
		response.ServerErrorResponse(w, s.Logger, err)
	}
}

func (s *server) handleListPosts(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		Tags  []string
		store.Filters
	}

	qs := r.URL.Query()

	input.Title = request.ReadString(qs, "title", "")
	input.Tags = request.ReadCSV(qs, "tags", []string{})

	input.Filters.Page = request.ReadInt(qs, "page", 1)
	input.Filters.Limit = request.ReadInt(qs, "limit", 20)

	input.Filters.Sort = request.ReadString(qs, "sort", "id")

	errs := request.ValidateInput(input.Filters)
	if errs != nil {
		response.FailedValidationResponse(w, errs)
		return
	}

	repo := store.NewPostRepository(s.DB)

	posts, err := repo.GetAll(input.Title, input.Tags, input.Filters)
	if err != nil {
		response.ServerErrorResponse(w, s.Logger, err)
		return
	}

	if err := response.JSONResponse(w, http.StatusOK, response.Envelope{"posts": posts}); err != nil {
		response.ServerErrorResponse(w, s.Logger, err)
	}
}
