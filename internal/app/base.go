package app

import (
	"context"
	"database/sql"
	"flag"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	Port  int
	Env   string
	Dsn string
}

type server struct {
	DB     *sql.DB
	Router *mux.Router
	Logger *logrus.Logger
	Config config
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

	s.getConfig()

	s.routes()

	s.Logger.Info("connecting the database")
	db, err := openDB(s.Config)
	if err != nil {
		s.Logger.WithError(err).Fatal("an error occurred while connecting the database")
	}

	s.DB = db
}

func (s *server) getConfig() {
	var cfg config

	s.Logger.Info("we are getting env values")
	err := godotenv.Load()
	if err != nil {
		s.Logger.WithError(err).Fatal("something went wrong while getting env")
	}

	flag.IntVar(&cfg.Port, "port", 3000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.Dsn, "db-dsn", os.Getenv("DB_URI"), "PostgreSQL DSN")
	flag.Parse()

	s.Config = cfg
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Dsn)
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
