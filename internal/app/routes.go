package app

import (
	"github.com/gorilla/mux"
	"github.com/nebisin/api_structure/pkg/response"
	"net/http"
)

func (s *server) routes() {
	s.Router = mux.NewRouter()

	s.Router.NotFoundHandler = http.HandlerFunc(response.NotFoundResponse)
	s.Router.MethodNotAllowedHandler = http.HandlerFunc(response.MethodNotAllowedResponse)

	s.Router.HandleFunc("/v1/healthcheck", s.handleHealthCheck)
	s.Router.HandleFunc("/v1/posts", s.handleCreatePost).Methods(http.MethodPost)
	s.Router.HandleFunc("/v1/posts/{id}", s.handleShowPost).Methods(http.MethodGet)
	s.Router.HandleFunc("/v1/posts", s.handleListPosts).Methods(http.MethodGet)
	s.Router.HandleFunc("/v1/posts/{id}", s.handleUpdatePost).Methods(http.MethodPatch)
	s.Router.HandleFunc("/v1/posts/{id}", s.handleDeletePost).Methods(http.MethodDelete)
}

func (s *server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	err := response.JSONResponse(w, http.StatusOK, map[string]bool{"ok": true})
	if err != nil {
		response.ServerErrorResponse(w, s.Logger, err)
	}
}
