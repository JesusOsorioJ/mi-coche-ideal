package vehicles_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"mi-coche-ideal/internal/testutils"
	"mi-coche-ideal/internal/vehicles"
)

func TestVehicleEndpoints(t *testing.T) {
	var resp *httptest.ResponseRecorder 


	gin.SetMode(gin.TestMode)
	DB := testutils.SetupDB(t)
	router := gin.Default()

	router.POST("/auth/signup", testutils.SignupHandler(DB))
	router.POST("/auth/login", testutils.LoginHandler(DB))

	token := testutils.SignupAndLogin(t, router, "vendedor@example.com", "vehiculo123")

	vehicleRepo := vehicles.NewVehicleRepository(DB)
	vehicleService := vehicles.NewVehicleService(vehicleRepo)
	vehicleHandler := vehicles.NewVehicleHandler(vehicleService)

	api := router.Group("/api")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
	})
	vehicleHandler.RegisterRoutes(api)

	vehicle := vehicles.Vehicle{
		Brand:      "Toyota",
		Model:      "Corolla",
		Year:       2020,
		Price:      15000,
		Kilometers: 30000,
		MainPhoto:  "https://example.com/car.jpg",
	}
	vJSON, _ := json.Marshal(vehicle)
	req := httptest.NewRequest(http.MethodPost, "/api/vehicles/", bytes.NewBuffer(vJSON))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code)

	req = httptest.NewRequest(http.MethodGet, "/api/vehicles/?brand=Toyota", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
}
