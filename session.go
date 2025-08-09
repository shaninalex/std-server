package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/uptrace/bun"
)

var _store *CookieStoreDatabase

func GetStore() *CookieStoreDatabase {
	return _store
}

func NewCookieStoreDatabase(db *bun.DB, keyPairs ...[]byte) *CookieStoreDatabase {
	st := &CookieStoreDatabase{
		db:     db,
		codecs: securecookie.CodecsFromPairs(keyPairs...),
		options: &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7, // 7 days
			HttpOnly: true,
		},
	}
	_store = st
	return st
}

type CookieStoreDatabase struct {
	db      *bun.DB
	codecs  []securecookie.Codec
	options *sessions.Options
}

func (s *CookieStoreDatabase) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *CookieStoreDatabase) MaxAge(age int) {
	s.options.MaxAge = age
	for _, c := range s.codecs {
		if sc, ok := c.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

func (s *CookieStoreDatabase) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	session.Options = &(*s.options)
	session.IsNew = true

	cookie, err := r.Cookie(name)
	if err != nil {
		// no cookie = new session
		return session, nil
	}

	// decode session ID from cookie
	var sessionID string
	if err := securecookie.DecodeMulti(name, cookie.Value, &sessionID, s.codecs...); err != nil {
		return session, nil // invalid cookie → treat as new
	}

	// load from DB
	ctx := r.Context()
	var sm SessionModel
	err = s.db.NewSelect().Model(&sm).Where("id = ?", sessionID).Scan(ctx)
	if err != nil {
		return session, nil // not found → new session
	}

	// decode Values
	buf := bytes.NewBuffer(sm.Data)
	dec := gob.NewDecoder(buf)
	values := make(map[interface{}]interface{})
	if err := dec.Decode(&values); err == nil {
		session.Values = values
	}
	session.ID = sessionID
	session.IsNew = false
	return session, nil
}

func (s *CookieStoreDatabase) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// generate ID if new
	if session.ID == "" {
		session.ID = uuid.NewString()
	}

	// --- make sure user_id is string, not uuid.UUID ---
	if v, ok := session.Values["user_id"].(uuid.UUID); ok {
		session.Values["user_id"] = v.String()
	}

	// serialize Values to []byte
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(session.Values); err != nil {
		return err
	}

	// parse session.ID -> uuid.UUID for DB
	uid, err := uuid.Parse(session.ID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	// parse user_id string -> uuid.UUID for DB
	var userID uuid.UUID
	if v, ok := session.Values["user_id"].(string); ok && v != "" {
		if parsed, err := uuid.Parse(v); err == nil {
			userID = parsed
		}
	}

	// upsert into DB
	ctx := r.Context()
	sm := &SessionModel{
		ID:        uid,
		UserID:    userID,
		Data:      buf.Bytes(),
		ExpiresAt: time.Now().Add(time.Duration(session.Options.MaxAge) * time.Second),
	}

	_, err = s.db.NewInsert().
		Model(sm).
		On("CONFLICT (id) DO UPDATE").
		Set("data = EXCLUDED.data, expires_at = EXCLUDED.expires_at").
		Exec(ctx)
	if err != nil {
		return err
	}

	// set cookie with signed session ID
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}
