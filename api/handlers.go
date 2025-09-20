package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shaninalex/std-server/pkg"
	"golang.org/x/crypto/bcrypt"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_user, err := pkg.GetUserByField(r.Context(), pkg.GetDB(), "email", req.Email)
	if err == nil && _user != nil {
		http.Error(w, "Email already taken", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	user := &pkg.UserModel{
		ID:           uuid.NewString(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = pkg.SaveUser(r.Context(), pkg.GetDB(), user)
	if err != nil {
		http.Error(w, "Could not save user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := pkg.GetUserByField(r.Context(), pkg.GetDB(), "email", req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, _ := pkg.GetStore().Get(r, "app_session")
	session.Values["user_id"] = user.ID
	session.Values["user_email"] = user.Email
	session.Options.MaxAge = 86400 * 7 // 7 днів

	if err := pkg.GetStore().Save(r, w, session); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save session: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged_in"})
}
