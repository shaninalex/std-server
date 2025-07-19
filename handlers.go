package main

import (
	"fmt"
	"net/http"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func HandleProtected(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(ContextUserIDKey)
	responseText := fmt.Sprintf("Profile ID: %s", userId)
	w.Write([]byte(responseText))
}
