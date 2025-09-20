package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Ping(t *testing.T) {
	router := NewRouter()
	router.GET("/ping", HandlePing)
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/ping")
	require.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", string(body))
}
