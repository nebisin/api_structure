package store

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Post struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags,omitempty"`
	Version   int32     `json:"version"`
}
