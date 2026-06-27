package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"vehicles/internal/db/storage"
	"vehicles/internal/models"
)

var validTypes = map[string]bool{
	"Car":   true,
	"Bike":  true,
	"Truck": true,
	"SUV":   true,
	"Van":   true,
	"Bus":   true,
	"Other": true,
}

var validCategories = map[string]bool{
	"Normal":   true,
	"Moderate": true,
	"Premium":  true,
}

type Service struct {
	VehicleRepo storage.Storage
}

func NewService(repo storage.Storage) *Service {
	return &Service{
		VehicleRepo: repo,
	}
}

func (s *Service) CreateVehicle(vehicles *models.Vehicle) error {

	if err := s.ValidateVehicle(vehicles); err != nil {
		return fmt.Errorf("failed to validate vehicle : %w", err)
	}

	if err := s.CheckDuplicateVehicle(vehicles); err != nil {
		return fmt.Errorf("failed to check duplicate vehicle : %w", err)
	}

	vehicles.CreatedAt = time.Now().UTC()

	err := s.VehicleRepo.CreateVehicle(vehicles)
	if err != nil {
		return fmt.Errorf("failed to create vehicle : %w", err)
	}
	return nil
}

func (s *Service) UpdateVehicle(vehicle *models.Vehicle) error {
	if err := s.ValidateVehicle(vehicle); err != nil {
		return fmt.Errorf("failed to validate vehicle: %w", err)
	}

	if err := s.CheckDuplicateForUpdate(vehicle); err != nil {
		return fmt.Errorf("failed to check vehicle: %w", err)
	}

	err := s.VehicleRepo.UpdateVehicle(vehicle)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}
	return nil
}

func (s *Service) DeleteVehicle(id int64) error {
	if err := s.VehicleRepo.DeleteVehicle(id); err != nil {
		return fmt.Errorf("failed to delete vehicle : %w", err)
	}

	return nil
}

func (s *Service) GetVehicleById(id int64) (*models.Vehicle, error) {
	vehicle, err := s.VehicleRepo.GetVehicleById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle by id : %w", err)
	}
	return vehicle, nil
}

func (s *Service) GetAllVehicles(page, limit int) (*models.PaginationResponse, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	vehicles, total, err := s.VehicleRepo.GetAllVehicles(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicles: %w", err)
	}

	totalPages := (total + limit - 1) / limit

	response := &models.PaginationResponse{
		Page:        int64(page),
		Limit:       int64(limit),
		Total:       int64(total),
		TotalPages:  int64(totalPages),
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
		Data:        vehicles,
	}

	return response, nil
}

func (s *Service) ValidateVehicle(vehicle *models.Vehicle) error {
	if strings.TrimSpace(vehicle.Name) == "" {
		return errors.New("vehicle name is required")
	}

	if strings.TrimSpace(vehicle.Model) == "" {
		return errors.New("vehicle model is required")
	}

	if !validTypes[vehicle.Type] {
		return errors.New("invalid vehicle type")
	}

	if !validCategories[vehicle.Category] {
		return errors.New("invalid vehicle category")
	}

	return nil
}

func (s *Service) CheckDuplicateVehicle(vehicle *models.Vehicle) error {
	exist, err := s.VehicleRepo.ExistByName(vehicle.Name)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("vehicle already exist")
	}

	return nil
}

func (s *Service) CheckDuplicateForUpdate(vehicle *models.Vehicle) error {
	exist, err := s.VehicleRepo.ExistsByNameExceptId(vehicle.Name, vehicle.Id)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("vehicle already exist")
	}

	return nil
}
