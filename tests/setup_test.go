package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mi-coche-ideal/internal/auth"
	"mi-coche-ideal/internal/orders"
	"mi-coche-ideal/internal/vehicles"
)

var DB *gorm.DB

func TestMain(m *testing.M) {
	_ = godotenv.Load("../.env")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a la base de datos: %v", err)
	}

	DB = db
	DB.Exec("TRUNCATE users, vehicles, orders RESTART IDENTITY CASCADE")
	_ = DB.AutoMigrate(&auth.User{}, &vehicles.Vehicle{}, &orders.Order{})

	code := m.Run()
	os.Exit(code)
}
