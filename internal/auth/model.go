package auth

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email" gorm:"unique"`
	Password  string    `json:"password,omitempty" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
