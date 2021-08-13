package store

import (
	"errors"
	"github.com/nebisin/api_structure/pkg/auth"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Post struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags,omitempty"`
	Version   int32     `json:"version"`
}

type User struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name string `json:"name" validate:"required,max=500"`
	Email string `json:"email" validate:"required,email"`
	Password auth.Password `json:"-"`
	Activated bool `json:"activated"`
	Version int `json:"-"`
}