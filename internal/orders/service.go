package orders

import (
	"errors"
)

type OrderService struct {
	repo *OrderRepository
}

func NewOrderService(repo *OrderRepository) *OrderService {
	return &OrderService{repo}
}

func (s *OrderService) Create(order *Order) error {
	order.Status = StatusPending
	return s.repo.Create(order)
}

func (s *OrderService) GetByID(id uint) (*Order, error) {
	return s.repo.GetByID(id)
}

func (s *OrderService) GetAllByUser(userID uint) ([]Order, error) {
	return s.repo.GetAllByUser(userID)
}

func (s *OrderService) UpdateStatus(order *Order, newStatus OrderStatus) error {
	validTransitions := map[OrderStatus][]OrderStatus{
		StatusPending:   {StatusPaid},
		StatusPaid:      {StatusDelivered},
		StatusDelivered: {},
	}

	allowed := false
	for _, s := range validTransitions[order.Status] {
		if s == newStatus {
			allowed = true
			break
		}
	}

	if !allowed {
		return errors.New("invalid status transition")
	}

	order.Status = newStatus
	return s.repo.Update(order)
}
