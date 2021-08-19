package app

import (
	"context"
	"github.com/nebisin/api_structure/internal/store"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func (s *server) contextSetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (s *server) contextGetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}