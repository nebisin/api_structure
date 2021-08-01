package store

import "database/sql"

type postRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{
		DB: db,
	}
}

func (r *postRepository) Insert(post *Post) error {
	return nil
}

func (r *postRepository) Get(id int64) (*Post, error) {
	return nil, nil
}

func (r *postRepository) Update(post *Post) error {
	return nil
}

func (r *postRepository) Delete(id int64) error {
	return nil
}
