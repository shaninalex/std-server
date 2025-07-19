package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Ping(t *testing.T) {
	router := NewBackendRouter()
	router.UnprotectedHandle("/ping", HandlePing)
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/ping")
	require.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", string(body))
}
