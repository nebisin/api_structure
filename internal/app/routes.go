package app

import (
	"github.com/gorilla/mux"
	"github.com/nebisin/api_structure/pkg/response"
	"net/http"
)

func (s *server) routes() {
	s.logger.Info("initializing the routes")

	s.router = mux.NewRouter()

	s.router.Use(s.recoverPanic)
	s.router.Use(s.enableCORS)
	s.router.Use(s.rateLimit)
	s.router.Use(s.authenticate)

	s.router.NotFoundHandler = http.HandlerFunc(response.NotFoundResponse)
	s.router.MethodNotAllowedHandler = http.HandlerFunc(response.MethodNotAllowedResponse)

	apiV1 := s.router.PathPrefix("/api/v1").Subrouter()

	apiV1.HandleFunc("/healthcheck", s.handleHealthCheck)

	apiV1.HandleFunc("/posts", s.requirePermission("posts:write", s.handleCreatePost)).Methods(http.MethodPost)
	apiV1.HandleFunc("/posts/{id}", s.requirePermission("posts:read", s.handleShowPost)).Methods(http.MethodGet)
	apiV1.HandleFunc("/posts", s.requirePermission("posts:read", s.handleListPosts)).Methods(http.MethodGet)
	apiV1.HandleFunc("/posts/{id}", s.requirePermission("posts:write", s.handleUpdatePost)).Methods(http.MethodPatch)
	apiV1.HandleFunc("/posts/{id}", s.requirePermission("posts:write", s.handleDeletePost)).Methods(http.MethodDelete)

	apiV1.HandleFunc("/users", s.handleRegisterUser).Methods(http.MethodPost)
	apiV1.HandleFunc("/users/activated", s.handleActivateUser).Methods(http.MethodPut)

	apiV1.HandleFunc("/tokens/authentication", s.handleCreateAuthenticationToken).Methods(http.MethodPost)
}

func (s *server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	err := response.JSONResponse(w, http.StatusOK, response.Envelope{
		"status":      "available",
		"environment": s.config.env,
		"version":     version,
	})
	if err != nil {
		response.ServerErrorResponse(w, r, s.logger, err)
	}
}
