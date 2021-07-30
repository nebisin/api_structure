package app

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
)

type server struct {
	DB     *sql.DB
	Router *mux.Router
	Logger *log.Logger
}

func NewServer() *server {
	return &server{}
}

func (s *server) Run() {
	s.routes()
}
