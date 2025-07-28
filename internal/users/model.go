package users

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"` // omitir en JSON
	CreatedAt time.Time
	UpdatedAt time.Time
}
