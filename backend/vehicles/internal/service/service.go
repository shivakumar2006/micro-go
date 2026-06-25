package service

import (
	"fmt"
	"vehicles/internal/models"
	"vehicles/internal/repository"
)

type Service struct {
	VehicleRepo *repository.VehicleRepository
}

func NewService(vr *repository.VehicleRepository) *Service {
	return &Service{
		VehicleRepo: vr,
	}
}

func (s *Service) CreateVehicle(vehicles *models.Vehicle) error {
	err := s.VehicleRepo.CreateVehicle(vehicles)
	if err != nil {
		return fmt.Errorf("failed to create vehicle : %w", err)
	}
	return nil
}
