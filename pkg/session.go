package pkg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var _store *CookieStoreDatabase

func GetStore() *CookieStoreDatabase {
	return _store
}

func NewCookieStoreDatabase(db *sql.DB, keyPairs ...[]byte) *CookieStoreDatabase {
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
	db      *sql.DB
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
		return session, nil
	}

	var sessionID string
	if err := securecookie.DecodeMulti(name, cookie.Value, &sessionID, s.codecs...); err != nil {
		return session, nil
	}

	uid, err := uuid.Parse(sessionID)
	if err != nil {
		return session, nil
	}

	sm, err := GetSessionByID(r.Context(), s.db, uid.String())
	if err != nil {
		return session, nil
	}

	if sm.Data != "" {
		values := make(map[string]interface{})
		if err := json.Unmarshal([]byte(sm.Data), &values); err == nil {
			for k, v := range values {
				session.Values[k] = v
			}
		} else {
			fmt.Println("failed to unmarshal session data:", err)
		}
	}
	session.ID = uid.String()
	session.IsNew = false
	return session, nil
}

func (s *CookieStoreDatabase) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.ID == "" {
		session.ID = uuid.NewString()
	}

	// normalize Values into map[string]interface{}
	values := make(map[string]interface{}, len(session.Values))
	for k, v := range session.Values {
		if ks, ok := k.(string); ok {
			values[ks] = v
		}
	}

	// serialize as JSON
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}

	// parse session.ID
	uid, err := uuid.Parse(session.ID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	// parse user_id
	var userID string
	if v, ok := values["user_id"].(string); ok {
		userID = v
	}

	ctx := r.Context()
	sm := &SessionModel{
		ID:        uid.String(),
		UserID:    userID,
		Data:      string(data),
		ExpiresAt: time.Now().Add(time.Duration(session.Options.MaxAge) * time.Second),
		CreatedAt: time.Now(),
	}

	err = SaveSession(ctx, s.db, sm)
	if err != nil {
		return err
	}

	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}
