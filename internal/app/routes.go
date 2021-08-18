package app

import (
	"github.com/gorilla/mux"
	"github.com/nebisin/api_structure/pkg/response"
	"net/http"
)

func (s *server) routes() {
	s.logger.Info("initializing the routes")

	s.router = mux.NewRouter()

	s.router.Use(s.rateLimit)
	s.router.Use(s.recoverPanic)

	s.router.NotFoundHandler = http.HandlerFunc(response.NotFoundResponse)
	s.router.MethodNotAllowedHandler = http.HandlerFunc(response.MethodNotAllowedResponse)

	s.router.HandleFunc("/v1/healthcheck", s.handleHealthCheck)
	s.router.HandleFunc("/v1/posts", s.handleCreatePost).Methods(http.MethodPost)
	s.router.HandleFunc("/v1/posts/{id}", s.handleShowPost).Methods(http.MethodGet)
	s.router.HandleFunc("/v1/posts", s.handleListPosts).Methods(http.MethodGet)
	s.router.HandleFunc("/v1/posts/{id}", s.handleUpdatePost).Methods(http.MethodPatch)
	s.router.HandleFunc("/v1/posts/{id}", s.handleDeletePost).Methods(http.MethodDelete)

	s.router.HandleFunc("/v1/users", s.handleRegisterUser).Methods(http.MethodPost)
	s.router.HandleFunc("/v1/users/activated", s.handleActivateUser).Methods(http.MethodPut)
}

func (s *server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	err := response.JSONResponse(w, http.StatusOK, response.Envelope{
		"status": "available",
		"environment": s.config.env,
		"version": version,
	})
	if err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
	}
}
