package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func SignupAndLogin(t *testing.T, router http.Handler, email, password string) string {
	var resp *httptest.ResponseRecorder 

	
	body := map[string]string{
		"name":     "Test User",
		"email":    email,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	req = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, 200, resp.Code)

	var result map[string]string
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))

	return result["access_token"]
}
