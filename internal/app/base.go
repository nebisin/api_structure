package app

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	dsn  string
}

type server struct {
	db     *sql.DB
	router *mux.Router
	logger *logrus.Logger
	config config
}

func NewServer() *server {
	return &server{}
}

func (s *server) Run() {
	s.logger = logrus.New()
	s.logger.SetOutput(os.Stdout)
	s.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	s.getConfig()

	s.routes()

	s.logger.Info("connecting the database")
	db, err := openDB(s.config)
	if err != nil {
		s.logger.WithError(err).Fatal("an error occurred while connecting the database")
	}
	s.db = db

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.port),
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	s.logger.WithField("port", srv.Addr).Info("starting the server")
	s.logger.Fatal(srv.ListenAndServe())
}

func (s *server) getConfig() {
	var cfg config

	s.logger.Info("we are getting env values")
	err := godotenv.Load()
	if err != nil {
		s.logger.WithError(err).Fatal("something went wrong while getting env")
	}

	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("DB_URI"), "PostgreSQL DSN")
	flag.Parse()

	s.config = cfg
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)
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
