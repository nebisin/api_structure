package app

import (
	"context"
	"database/sql"
	"flag"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nebisin/api_structure/internal/mailer"
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
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
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
	mailer mailer.Mailer
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

	s.mailer = mailer.New(s.config.smtp.host, s.config.smtp.port, s.config.smtp.username, s.config.smtp.password, s.config.smtp.sender)

	s.logger.Info("connecting the database")
	db, err := openDB(s.config)
	if err != nil {
		s.logger.WithError(err).Fatal("an error occurred while connecting the database")
	}
	defer db.Close()
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

	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")

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
