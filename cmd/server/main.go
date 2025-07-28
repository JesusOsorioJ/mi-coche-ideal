package main

import (
	"os"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"mi-coche-ideal/internal/auth"
	"mi-coche-ideal/internal/csv"
	"mi-coche-ideal/internal/database"
	"mi-coche-ideal/internal/logging"
	"mi-coche-ideal/internal/metrics"
	"mi-coche-ideal/internal/middleware"
	"mi-coche-ideal/internal/orders"
	"mi-coche-ideal/internal/vehicles"
)

func main() {
	// Cargar .env
	_ = godotenv.Load()

	err := godotenv.Load()
if err != nil {
	log.Println("No .env file found")
}

log.Println("JWT_SECRET cargado:", os.Getenv("JWT_SECRET"))

	// Inicializar logger
	logging.InitLogger()
	log := logging.Log
	log.Info().Msg("🚀 Iniciando microservicio Mi Coche Ideal...")

	// Inicializar DB
	database.InitDB()
	db := database.DB

	// Iniciar rutina concurrente
	go csv.StartPriceUpdater(db, "price_updates.csv")

	// Inicializar router y métricas
	router := gin.Default()
	metrics.Init()
	router.Use(metrics.PrometheusMiddleware())

	// AUTH
	authService := auth.NewAuthService(db)
	authHandler := auth.NewAuthHandler(authService)

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/signup", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// RUTAS PROTEGIDAS
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	// VEHICLES
	vehicleRepo := vehicles.NewVehicleRepository(db)
	vehicleService := vehicles.NewVehicleService(vehicleRepo)
	vehicleHandler := vehicles.NewVehicleHandler(vehicleService)
	vehicleHandler.RegisterRoutes(protected)

	// ORDERS
	orderRepo := orders.NewOrderRepository(db)
	orderService := orders.NewOrderService(orderRepo)
	orderHandler := orders.NewOrderHandler(orderService)
	orderHandler.RegisterRoutes(protected)

	// Rutas públicas
	
	router.GET("/metrics", metrics.MetricsHandler())
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Levantar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Info().Str("port", port).Msg("Servidor corriendo en puerto")

	router.Run(":" + port)
}
