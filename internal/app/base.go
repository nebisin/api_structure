package app

import (
	"context"
	"database/sql"
	"flag"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"os"
	"sync"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	dsn  string
}

type server struct {
	db      *sql.DB
	router  *mux.Router
	logger  *logrus.Logger
	config  config
	limiter struct {
		mu      sync.Mutex
		clients map[string]*client
	}
}

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
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

	s.setupLimiter()

	s.logger.Info("connecting the database")
	db, err := openDB(s.config)
	if err != nil {
		s.logger.WithError(err).Fatal("an error occurred while connecting the database")
	}
	s.db = db

	if err := s.serve(); err != nil {
		s.logger.WithError(err).Fatal("an error occurred while starting the server")
	}
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

func (s *server) setupLimiter() {
	s.limiter.clients = make(map[string]*client)

	go func() {
		for {
			time.Sleep(time.Minute)

			s.limiter.mu.Lock()

			for ip, client := range s.limiter.clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(s.limiter.clients, ip)
				}
			}

			s.limiter.mu.Unlock()
		}
	}()
}
