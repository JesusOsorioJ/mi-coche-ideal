package orders

import (
	"time"

	"mi-coche-ideal/internal/users"
	"mi-coche-ideal/internal/vehicles"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pendiente"
	StatusPaid      OrderStatus = "pagada"
	StatusDelivered OrderStatus = "entregada"
)

type Order struct {
	ID         uint          `gorm:"primaryKey"`
	UserID     uint          `json:"user_id"`
	User       users.User    `gorm:"foreignKey:UserID"`
	VehicleID  uint          `json:"vehicle_id"`
	Vehicle    vehicles.Vehicle `gorm:"foreignKey:VehicleID"`
	Status     OrderStatus   `json:"status"`
	TotalPrice float64       `json:"total_price"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
