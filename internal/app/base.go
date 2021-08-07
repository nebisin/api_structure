package app

import (
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type server struct {
	DB     *sql.DB
	Router *mux.Router
	Logger *logrus.Logger
}

func NewServer() *server {
	return &server{}
}

func (s *server) Run() {
	s.Logger = logrus.New()
	s.Logger.SetOutput(os.Stdout)
	s.Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	s.Logger.Info("we are getting env values")
	err := godotenv.Load()
	if err != nil {
		s.Logger.WithError(err).Fatal("something went wrong while getting env")
	}

	s.routes()

	s.Logger.Info("connecting the database")
	db, err := openDB()
	if err != nil {
		s.Logger.WithError(err).Fatal("an error occurred while connecting the database")
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
