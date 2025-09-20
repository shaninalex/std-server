package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// SaveUser inserts a new user into the database
func SaveUser(ctx context.Context, db *sql.DB, user *UserModel) error {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	query := `
	INSERT INTO users (id, name, email, password_hash, active, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

// GetUserByField fetches a user by any column (e.g., "id" or "email").
func GetUserByField(ctx context.Context, db *sql.DB, field, value string) (*UserModel, error) {
	user := &UserModel{}

	query := fmt.Sprintf(`
	SELECT id, name, email, password_hash, active, created_at, updated_at
	FROM users
	WHERE %s = ?
	LIMIT 1
	`, field)

	row := db.QueryRowContext(ctx, query, value)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return user, nil
}

// GetSessionByID loads a session from DB by session ID.
func GetSessionByID(ctx context.Context, db *sql.DB, id string) (*SessionModel, error) {
	sm := &SessionModel{}

	query := `
	SELECT id, user_id, data, expires_at, created_at
	FROM user_sessions
	WHERE id = ?
	LIMIT 1
	`

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&sm.ID,
		&sm.UserID,
		&sm.Data,
		&sm.ExpiresAt,
		&sm.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return sm, nil
}

// SaveSession inserts or updates a session in the DB.
func SaveSession(ctx context.Context, db *sql.DB, sm *SessionModel) error {
	if sm.CreatedAt.IsZero() {
		sm.CreatedAt = time.Now()
	}
	if sm.ExpiresAt.IsZero() {
		sm.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // default 7 days
	}

	query := `
	INSERT INTO user_sessions (id, user_id, data, expires_at, created_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		data = excluded.data,
		expires_at = excluded.expires_at
	`

	_, err := db.ExecContext(ctx, query,
		sm.ID,
		sm.UserID,
		sm.Data,
		sm.ExpiresAt,
		sm.CreatedAt,
	)
	return err
}
