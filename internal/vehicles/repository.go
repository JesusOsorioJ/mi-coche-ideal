package vehicles

import (
	"gorm.io/gorm"
)

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db}
}

func (r *VehicleRepository) Create(vehicle *Vehicle) error {
	return r.db.Create(vehicle).Error
}

func (r *VehicleRepository) GetByID(id uint) (*Vehicle, error) {
	var vehicle Vehicle
	err := r.db.First(&vehicle, id).Error
	return &vehicle, err
}

func (r *VehicleRepository) Update(vehicle *Vehicle) error {
	return r.db.Save(vehicle).Error
}

func (r *VehicleRepository) Delete(id uint) error {
	return r.db.Delete(&Vehicle{}, id).Error
}

func (r *VehicleRepository) ListFiltered(brand, model string, year int, minPrice, maxPrice float64, page, limit int) ([]Vehicle, error) {
	var vehicles []Vehicle
	query := r.db.Model(&Vehicle{})

	if brand != "" {
		query = query.Where("brand ILIKE ?", "%"+brand+"%")
	}
	if model != "" {
		query = query.Where("model ILIKE ?", "%"+model+"%")
	}
	if year > 0 {
		query = query.Where("year = ?", year)
	}
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&vehicles).Error
	return vehicles, err
}
