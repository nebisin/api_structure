package app

import (
	"errors"
	"fmt"
	"github.com/nebisin/api_structure/internal/store"
	"github.com/nebisin/api_structure/pkg/response"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"strings"
	"time"
)

func (s *server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				response.ServerErrorResponse(w, r, s.logger, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (s *server) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			response.ServerErrorResponse(w, r, s.logger, err)
		}

		s.limiter.mu.Lock()

		if _, found := s.limiter.clients[ip]; !found {
			s.limiter.clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		s.limiter.clients[ip].lastSeen = time.Now()

		if !s.limiter.clients[ip].limiter.Allow() {
			s.limiter.mu.Unlock()
			response.RateLimitExceededResponse(w, r)
			return
		}

		s.limiter.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func (s *server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = s.contextSetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			response.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		if len(token) > 26 {
			response.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		userRepo := store.NewUserRepository(s.db)
		user, err := userRepo.GetForToken(store.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				response.InvalidAuthenticationTokenResponse(w, r)
			default:
				response.ServerErrorResponse(w, r, s.logger, err)
			}
			return
		}

		r = s.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}