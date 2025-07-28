
package orders_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strconv"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "mi-coche-ideal/internal/orders"
    "mi-coche-ideal/internal/testutils"
    "mi-coche-ideal/internal/vehicles"
)

func TestOrderFlow(t *testing.T) {
	var resp *httptest.ResponseRecorder 

	
    gin.SetMode(gin.TestMode)
    DB := testutils.SetupDB(t)
    router := gin.Default()

    router.POST("/auth/signup", testutils.SignupHandler(DB))
    router.POST("/auth/login", testutils.LoginHandler(DB))
    token := testutils.SignupAndLogin(t, router, "comprador@example.com", "orden123")
    _ = token

    vehicleRepo := vehicles.NewVehicleRepository(DB)
    vehicleService := vehicles.NewVehicleService(vehicleRepo)
    vehicleHandler := vehicles.NewVehicleHandler(vehicleService)

    orderRepo := orders.NewOrderRepository(DB)
    orderService := orders.NewOrderService(orderRepo)
    orderHandler := orders.NewOrderHandler(orderService)

    protected := router.Group("/api")
    protected.Use(func(c *gin.Context) {
        c.Set("user_id", uint(1))
    })
    vehicleHandler.RegisterRoutes(protected)
    orderHandler.RegisterRoutes(protected)

    vehicle := vehicles.Vehicle{
        Brand:      "Ford",
        Model:      "Focus",
        Year:       2018,
        Price:      12000,
        Kilometers: 50000,
        MainPhoto:  "https://example.com/focus.jpg",
    }
    vJSON, _ := json.Marshal(vehicle)
    req := httptest.NewRequest(http.MethodPost, "/api/vehicles/", bytes.NewBuffer(vJSON))
    req.Header.Set("Content-Type", "application/json")
    resp = httptest.NewRecorder()
    router.ServeHTTP(resp, req)
    assert.Equal(t, 201, resp.Code)

    var createdVehicle vehicles.Vehicle
    err := json.Unmarshal(resp.Body.Bytes(), &createdVehicle)
    require.NoError(t, err)

    orderReq := map[string]interface{}{
        "vehicle_id":  createdVehicle.ID,
        "total_price": createdVehicle.Price,
    }
    orderJSON, _ := json.Marshal(orderReq)
    req = httptest.NewRequest(http.MethodPost, "/api/orders/", bytes.NewBuffer(orderJSON))
    req.Header.Set("Content-Type", "application/json")
    resp = httptest.NewRecorder()
    router.ServeHTTP(resp, req)
    assert.Equal(t, 201, resp.Code)

    var createdOrder orders.Order
    err = json.Unmarshal(resp.Body.Bytes(), &createdOrder)
    require.NoError(t, err)
    assert.Equal(t, orders.StatusPending, createdOrder.Status)

    update := map[string]string{"status": "pagada"}
    updateJSON, _ := json.Marshal(update)
    req = httptest.NewRequest(http.MethodPut, "/api/orders/"+strconv.Itoa(int(createdOrder.ID))+"/status", bytes.NewBuffer(updateJSON))
    req.Header.Set("Content-Type", "application/json")
    resp = httptest.NewRecorder()
    router.ServeHTTP(resp, req)
    assert.Equal(t, 200, resp.Code)

    var updatedOrder orders.Order
    err = json.Unmarshal(resp.Body.Bytes(), &updatedOrder)
    require.NoError(t, err)
    assert.Equal(t, orders.StatusPaid, updatedOrder.Status)
}
