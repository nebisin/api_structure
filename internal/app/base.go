package app

import (
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
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
	err := godotenv.Load()
	if err != nil {
		s.Logger.Fatal("Error loading .env file")
	}

	s.routes()

	db, err := openDB()
	if err != nil {
		s.Logger.Fatal(err)
	}

	s.DB = db
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DB_URI"))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
