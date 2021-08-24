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

// rateLimit is the rate limiter middleware
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

// authenticate is the authentication middleware
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

		user, err := s.models.Users.GetForToken(store.ScopeAuthentication, token)
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

func (s *server) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.contextGetUser(r)

		if user.IsAnonymous() {
			response.AuthenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.contextGetUser(r)

		if !user.Activated {
			response.InactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	return s.requireAuthenticatedUser(fn)
}

func (s *server) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := s.contextGetUser(r)

		permissions, err := s.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			response.ServerErrorResponse(w, r, s.logger, err)
			return
		}

		if !permissions.Include(code) {
			response.NotPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	return s.requireActivatedUser(fn)
}

func (s *server) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" && len(s.config.cors.trustedOrigins) != 0 {
			for i := range s.config.cors.trustedOrigins {
				if origin == s.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}

		if r.Method == http.MethodOptions {
			response.MethodNotAllowedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
