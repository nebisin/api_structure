package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type postRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{
		DB: db,
	}
}

func (r *postRepository) Insert(post *Post) error {
	query := `INSERT INTO posts (title, body, tags) 
VALUES ($1, $2, $3) 
RETURNING id, created_at, version`

	args := []interface{}{post.Title, post.Body, pq.Array(post.Tags)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&post.ID, &post.CreatedAt, &post.Version)
}

func (r *postRepository) Get(id int64) (*Post, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, created_at, title, body, tags, version
FROM posts
WHERE id = $1`

	var post Post

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.Title,
		&post.Body,
		pq.Array(&post.Tags),
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (r *postRepository) Update(post *Post) error {
	query := `UPDATE posts SET title=$1, body=$2, tags=$3, version= version + 1
WHERE id=$4 AND version = $5
RETURNING version`

	args := []interface{}{
		post.Title,
		post.Body,
		pq.Array(post.Tags),
		post.ID,
		post.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r *postRepository) Delete(id int64) error {
	query := `DELETE FROM posts
WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
