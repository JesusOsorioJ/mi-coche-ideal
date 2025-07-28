package vehicles

import "time"

type Vehicle struct {
	ID          uint      `gorm:"primaryKey"`
	Brand       string    `json:"brand"`
	Model       string    `json:"model"`
	Year        int       `json:"year"`
	Price       float64   `json:"price"`
	Kilometers  float64   `json:"kilometers"`
	MainPhoto   string    `json:"main_photo"` // URL
	Description string    `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
