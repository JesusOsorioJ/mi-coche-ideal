package database

import (
	"fmt"
	"log"
	"os"
	"time" // ⬅️ Faltaba esta importación

	"mi-coche-ideal/internal/orders"
	"mi-coche-ideal/internal/users"
	"mi-coche-ideal/internal/vehicles"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var db *gorm.DB
	var err error

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Reintentando conexión a la base de datos (%d/%d)...", i+1, maxRetries)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ Error al conectar a la base de datos: %v", err)
	}

	// Auto migración
	err = db.AutoMigrate(&users.User{}, &vehicles.Vehicle{}, &orders.Order{})
	if err != nil {
		log.Fatal("Failed to migrate models:", err)
	}

	DB = db
}
