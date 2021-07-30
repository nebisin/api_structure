package app

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
)

type handler struct {
	DB     *sql.DB
	Router *mux.Router
	Logger *log.Logger
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Run() {
	h.routes()
}
