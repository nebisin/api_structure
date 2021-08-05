package store

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Post struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags,omitempty"`
	Version   int32     `json:"version"`
}

type Filters struct {
	Page     int    `validate:"gt=0" json:"page"`
	PageSize int    `validate:"gt=0,lt=100" json:"page_size"`
	Sort     string `json:"sort"`
}
