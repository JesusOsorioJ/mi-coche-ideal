package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"mi-coche-ideal/internal/testutils"
)

func TestSignupAndLogin(t *testing.T) {
	var resp *httptest.ResponseRecorder 


	gin.SetMode(gin.TestMode)
	DB := testutils.SetupDB(t)
	router := gin.Default()

	router.POST("/auth/signup", testutils.SignupHandler(DB))
	router.POST("/auth/login", testutils.LoginHandler(DB))

	// Registro exitoso
	user := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonUser, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code)

	// Login exitoso
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	var result map[string]string
	_ = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NotEmpty(t, result["access_token"])

	// Login con contraseña incorrecta
	badLogin := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpass",
	}
	badJSON, _ := json.Marshal(badLogin)

	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(badJSON))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 401, resp.Code)
}
