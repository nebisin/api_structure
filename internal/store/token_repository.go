package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))

	token.Hash = hash[:]

	return token, nil
}

type tokenRepository struct {
	DB *sql.DB
}

func NewTokenRepository(db *sql.DB) *tokenRepository {
	return &tokenRepository{
		DB: db,
	}
}

func (r *tokenRepository) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = r.Insert(token)
	return token, err
}

func (r *tokenRepository) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope)
VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, args...)
	return err
}

func (r *tokenRepository) DeleteAllForUser(scope string, userID int64) error {
	query := `DELETE FROM tokens
WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, scope, userID)
	return err
}
