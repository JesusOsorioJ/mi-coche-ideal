package orders

import (
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) Create(order *Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) GetByID(id uint) (*Order, error) {
	var order Order
	err := r.db.Preload("User").Preload("Vehicle").First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) GetAllByUser(userID uint) ([]Order, error) {
	var orders []Order
	err := r.db.Preload("Vehicle").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) Update(order *Order) error {
	return r.db.Save(order).Error
}
