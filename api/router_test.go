package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Ping(t *testing.T) {
	router := NewRouter()
	router.GET("/ping", HandlePing)
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/ping")
	if err != nil {
		t.Errorf("Should not contain errors. Got %v", err)
	}

	body, _ := io.ReadAll(resp.Body)

	if http.StatusOK != resp.StatusCode {
		t.Errorf("Invalid status code. got %q, wanted %q", resp.StatusCode, http.StatusOK)
	}
	if "pong" != string(body) {
		t.Errorf("Invalid response body. got %q, wanted %q", string(body), "pong")
	}
}
