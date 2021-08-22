package store

import (
	"context"
	"database/sql"
	"time"
)

// Permissions slice will hold the permission codes.
type Permissions []string

// Include method checks the Permissions slice whether it contains
// a specific permission code or not.
func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}

	return false
}

type permissionRepository struct {
	DB *sql.DB
}

// GetAllForUser method returns all permission codes for a specific user
// in a Permissions slice.
func (r *permissionRepository) GetAllForUser(userID int64) (Permissions, error) {
	query := `
SELECT permissions.code
FROM permissions
INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
INNER JOIN users ON users_permissions.user_id = user.id
WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var permissions Permissions

	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
