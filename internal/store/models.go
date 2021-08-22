package store

import "database/sql"

type Models struct {
	Permissions permissionRepository
	Posts       postRepository
	Tokens      tokenRepository
	Users       userRepository
}

func NewModels(db *sql.DB) Models {
	return Models{
		Permissions: permissionRepository{db},
		Posts:       postRepository{db},
		Tokens:      tokenRepository{db},
		Users:       userRepository{db},
	}
}
