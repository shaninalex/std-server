package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const SessionTTL = 30 * 24 * time.Hour

type SessionModel struct {
	bun.BaseModel `bun:"table:user_sessions,alias:s"`

	ID        uuid.UUID `bun:"type:uuid,pk"`
	UserID    uuid.UUID `bun:"type:uuid,notnull"`
	Data      []byte    `bun:"data"`
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

func CreateSession(userID uuid.UUID, values map[any]any) *SessionModel {
	// serialize values
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(values); err != nil {
		panic(err) // or return error if you want
	}

	return &SessionModel{
		ID:        uuid.New(),
		UserID:    userID,
		Data:      buf.Bytes(),
		ExpiresAt: time.Now().Add(SessionTTL),
		CreatedAt: time.Now(),
	}
}
