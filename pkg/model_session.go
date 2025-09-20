package pkg

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SessionModel struct {
	bun.BaseModel `bun:"table:user_sessions,alias:s"`

	ID        uuid.UUID `bun:"type:uuid,pk"`
	UserID    uuid.UUID `bun:"type:uuid,notnull"`
	Data      string    `bun:"data,type:text"`
	ExpiresAt time.Time `bun:",notnull"`
	CreatedAt time.Time `bun:",notnull"`
}

func (s *SessionModel) BeforeInsert(ctx context.Context) (context.Context, error) {
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	return ctx, nil
}

func (s *SessionModel) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
