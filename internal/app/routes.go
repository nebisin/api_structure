package app

import (
	"github.com/gorilla/mux"
	"github.com/nebisin/api_structure/pkg/response"
	"net/http"
)

func (h *handler) routes() {
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/healthcheck", h.handleHealthCheck)
	h.Router.HandleFunc("/posts", h.handleCreatePost).Methods(http.MethodPost)
	h.Router.HandleFunc("/posts/{id:[0-9]+}", h.handleShowPost).Methods(http.MethodGet)
}

func (h *handler) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	err := response.JSONResponse(w, http.StatusOK, map[string]bool{"ok": true})
	if err != nil {
		h.Logger.Println(err)
		response.ServerErrorResponse(w)
	}
}
