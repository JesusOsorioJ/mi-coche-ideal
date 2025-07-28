package vehicles

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	service *VehicleService
}

func NewVehicleHandler(service *VehicleService) *VehicleHandler {
	return &VehicleHandler{service}
}

func (h *VehicleHandler) RegisterRoutes(rg *gin.RouterGroup) {
	v := rg.Group("/vehicles")
	v.POST("/", h.Create)
	v.GET("/", h.List)
	v.GET("/:id", h.GetByID)
	v.PUT("/:id", h.Update)
	v.DELETE("/:id", h.Delete)
}

func (h *VehicleHandler) Create(c *gin.Context) {
	var vehicle Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := h.service.Create(&vehicle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vehicle"})
		return
	}
	c.JSON(http.StatusCreated, vehicle)
}

func (h *VehicleHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	vehicle, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}
	c.JSON(http.StatusOK, vehicle)
}

func (h *VehicleHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var vehicle Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	vehicle.ID = uint(id)
	if err := h.service.Update(&vehicle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, vehicle)
}

func (h *VehicleHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *VehicleHandler) List(c *gin.Context) {
	brand := c.Query("brand")
	model := c.Query("model")
	year, _ := strconv.Atoi(c.Query("year"))
	minPrice, _ := strconv.ParseFloat(c.Query("price_min"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("price_max"), 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	vehicles, err := h.service.ListFiltered(brand, model, year, minPrice, maxPrice, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicles"})
		return
	}
	c.JSON(http.StatusOK, vehicles)
}
