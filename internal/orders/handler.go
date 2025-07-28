package orders

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DTO para crear orden
type OrderCreateDTO struct {
	VehicleID  uint    `json:"vehicle_id"`
	TotalPrice float64 `json:"total_price"`
}

type OrderHandler struct {
	service *OrderService
}

func NewOrderHandler(service *OrderService) *OrderHandler {
	return &OrderHandler{service}
}

func (h *OrderHandler) RegisterRoutes(rg *gin.RouterGroup) {
	o := rg.Group("/orders")
	o.POST("/", h.Create)
	o.GET("/", h.ListUserOrders)
	o.GET("/:id", h.GetByID)
	o.PUT("/:id/status", h.UpdateStatus)
}

func (h *OrderHandler) Create(c *gin.Context) {
	var input OrderCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userID := userIDInterface.(uint)

	order := &Order{
		UserID:     userID,
		VehicleID:  input.VehicleID,
		TotalPrice: input.TotalPrice,
		Status:     StatusPending,
	}

	if err := h.service.Create(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la orden"})
		return
	}

	// Recargar orden con relaciones
	fullOrder, err := h.service.GetByID(order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar detalles de la orden"})
		return
	}

	c.JSON(http.StatusCreated, fullOrder)
}

func (h *OrderHandler) ListUserOrders(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	orders, err := h.service.GetAllByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las órdenes"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order, err := h.service.GetByID(uint(id))
	if err != nil || order.UserID != c.MustGet("user_id").(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Orden no encontrada o acceso denegado"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		Status OrderStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	order, err := h.service.GetByID(uint(id))
	if err != nil || order.UserID != c.MustGet("user_id").(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Orden no encontrada o acceso denegado"})
		return
	}

	if err := h.service.UpdateStatus(order, input.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
