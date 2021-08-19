package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"github.com/nebisin/api_structure/pkg/auth"
	"time"
)

type User struct {
	ID        int64         `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Password  auth.Password `json:"-"`
	Activated bool          `json:"activated"`
	Version   int           `json:"-"`
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) Insert(user *User) error {
	query := `INSERT INTO users (name, email, password_hash, activated)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version`

	args := []interface{}{user.Name, user.Email, user.Password.Hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*User, error) {
	query := `SELECT id, created_at, name, email, password_hash, activated, version
FROM users
WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *userRepository) Update(user *User) error {
	query := `UPDATE users
SET name = $1, email = $2, password_hash=$3, activated=$4, version=version+1
WHERE id = $5 AND version = $6
RETURNING version`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r *userRepository) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
FROM users
INNER JOIN tokens
ON users.id = tokens.user_id
WHERE tokens.hash = $1 AND tokens.scope = $2 AND tokens.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
