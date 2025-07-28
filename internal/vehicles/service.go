package vehicles

type VehicleService struct {
	repo *VehicleRepository
}

func NewVehicleService(repo *VehicleRepository) *VehicleService {
	return &VehicleService{repo}
}

func (s *VehicleService) Create(vehicle *Vehicle) error {
	return s.repo.Create(vehicle)
}

func (s *VehicleService) GetByID(id uint) (*Vehicle, error) {
	return s.repo.GetByID(id)
}

func (s *VehicleService) Update(vehicle *Vehicle) error {
	return s.repo.Update(vehicle)
}

func (s *VehicleService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *VehicleService) ListFiltered(brand, model string, year int, minPrice, maxPrice float64, page, limit int) ([]Vehicle, error) {
	return s.repo.ListFiltered(brand, model, year, minPrice, maxPrice, page, limit)
}
