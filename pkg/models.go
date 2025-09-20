package pkg

import (
	"time"
)

// SessionModel represents a session for a user
type SessionModel struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Data      string    `db:"data"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type UserModel struct {
	ID           string    `db:"id" json:"id"` // UUID as string
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Active       bool      `db:"active" json:"active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
