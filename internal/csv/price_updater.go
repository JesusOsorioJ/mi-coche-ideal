package csv

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"mi-coche-ideal/internal/vehicles"
)

type PriceUpdate struct {
	VehicleID uint
	NewPrice  float64
}

func StartPriceUpdater(db *gorm.DB, filePath string) {
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		log.Println("⏱️ Ejecutando actualización de precios desde CSV...")
		updates, err := readCSV(filePath)
		if err != nil {
			log.Println("❌ Error al leer CSV:", err)
			return
		}

		applyUpdatesConcurrently(db, updates)
	})
	c.Start()
}

func readCSV(path string) ([]PriceUpdate, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var updates []PriceUpdate
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Suponemos que CSV tiene encabezado: vehicle_id,new_price
	for _, row := range records[1:] {
		id, _ := strconv.Atoi(row[0])
		price, _ := strconv.ParseFloat(row[1], 64)
		updates = append(updates, PriceUpdate{
			VehicleID: uint(id),
			NewPrice:  price,
		})
	}

	return updates, nil
}

func applyUpdatesConcurrently(db *gorm.DB, updates []PriceUpdate) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, update := range updates {
		wg.Add(1)
		go func(update PriceUpdate) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()

			tx := db.Begin()
			if err := tx.Model(&vehicles.Vehicle{}).
				Where("id = ?", update.VehicleID).
				Update("price", update.NewPrice).Error; err != nil {
				tx.Rollback()
				log.Printf("❌ Error al actualizar vehículo %d: %v", update.VehicleID, err)
				return
			}
			tx.Commit()
			log.Printf("✅ Precio actualizado para vehículo %d → %.2f", update.VehicleID, update.NewPrice)
		}(update)
	}

	wg.Wait()
	log.Println("✅ Finalizó rutina de precios.")
}

