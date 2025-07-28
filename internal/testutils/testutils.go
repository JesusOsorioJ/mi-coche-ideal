package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mi-coche-ideal/internal/auth"
)

func SetupDB(t *testing.T) *gorm.DB {
	err := godotenv.Load("../../.env")
	require.NoError(t, err)

	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	err = db.Exec("TRUNCATE users, vehicles, orders RESTART IDENTITY CASCADE").Error
	require.NoError(t, err)

	return db
}

func SignupHandler(DB *gorm.DB) gin.HandlerFunc {
	authService := auth.NewAuthService(DB)
	handler := auth.NewAuthHandler(authService)
	return handler.Register
}

func LoginHandler(DB *gorm.DB) gin.HandlerFunc {
	authService := auth.NewAuthService(DB)
	handler := auth.NewAuthHandler(authService)
	return handler.Login
}

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
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	require.NoError(t, err)

	return result["access_token"]
}
