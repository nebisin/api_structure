package app

import (
	"fmt"
	"github.com/nebisin/api_structure/pkg/response"
	"net/http"
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
		if !s.limiter.Allow() {
			response.RateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}