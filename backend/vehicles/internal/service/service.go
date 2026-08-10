package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"
	"vehicles/internal/db/storage"
	"vehicles/internal/models"
	storageCache "vehicles/internal/redis/storagecache"
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
	VehicleRepo  storage.Storage
	VehicleCache storageCache.StorageCache
}

func NewService(repo storage.Storage, cache storageCache.StorageCache) *Service {
	return &Service{
		VehicleRepo:  repo,
		VehicleCache: cache,
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

	// redis invalidate cache
	if err := s.VehicleCache.DeletePatterns("vehicles:*"); err != nil {
		log.Printf("failed to invalidate vehicle cache: %v", err)
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

	key := fmt.Sprintf("vehicle:%d", vehicle.Id)

	// redis
	if err := s.VehicleCache.Delete(key); err != nil {
		log.Printf("failed to delete vehicle from cache: %v", err)
	}

	if err := s.VehicleCache.DeletePatterns("vehicles:*"); err != nil {
		log.Printf("failed to delete vehicle from cache: %v", err)
	}

	return nil
}

func (s *Service) DeleteVehicle(id int64) error {
	if err := s.VehicleRepo.DeleteVehicle(id); err != nil {
		return fmt.Errorf("failed to delete vehicle : %w", err)
	}

	key := fmt.Sprintf("vehicle:%d", id)

	// redis
	if err := s.VehicleCache.Delete(key); err != nil {
		log.Printf("failed to delete vehicle from cache: %v", err)
	}

	if err := s.VehicleCache.DeletePatterns("vehicles:*"); err != nil {
		log.Printf("failed to delete vehicle from cache: %v", err)
	}

	return nil
}

func (s *Service) GetVehicleById(id int64) (*models.Vehicle, error) {
	key := fmt.Sprintf("vehicle:%d", id)

	// redis
	var vehicle models.Vehicle
	if err := s.VehicleCache.GetJSON(key, &vehicle); err == nil {
		return &vehicle, nil
	}

	vehicleFromDB, err := s.VehicleRepo.GetVehicleById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle by id : %w", err)
	}

	// if not present in redis then store in redis
	if err := s.VehicleCache.SetJSON(key, vehicleFromDB, 10*time.Minute); err != nil {
		slog.Error("failed to store vehicle in cache", "error", err)
	}

	return vehicleFromDB, nil
}

func (s *Service) GetAllVehicles(query models.VehicleQuery) (*models.PaginationResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.Limit < 1 {
		query.Limit = 10
	}

	if query.Limit > 100 {
		query.Limit = 100
	}

	// search
	query.Search = strings.TrimSpace(query.Search)

	// default sorting
	if query.SortBy == "" {
		query.SortBy = "created_at"
	}

	if query.Order == "" {
		query.Order = "desc"
	}

	validSortField := map[string]bool{
		"name":       true,
		"model":      true,
		"created_at": true,
	}

	if !validSortField[query.SortBy] {
		return nil, errors.New("invalid sort field")
	}

	if query.Order != "asc" && query.Order != "desc" {
		return nil, errors.New("invalid sort order")
	}

	if query.Type != "" && !validTypes[query.Type] {
		return nil, errors.New("invalid vehicle type")
	}

	if query.Category != "" && !validCategories[query.Category] {
		return nil, errors.New("invalid category type")
	}

	// redis
	key := fmt.Sprintf("vehicles:page=%d:limit=%d:search=%s:sort=%s:order=%s:type=%s:category=%s", query.Page, query.Limit, query.Search, query.SortBy, query.Order, query.Type, query.Category)

	var cached models.PaginationResponse
	if err := s.VehicleCache.GetJSON(key, &cached); err == nil {
		return &cached, nil
	}

	vehicles, total, err := s.VehicleRepo.GetAllVehicles(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicles: %w", err)
	}

	totalPages := (total + query.Limit - 1) / query.Limit

	response := &models.PaginationResponse{
		Page:        int64(query.Page),
		Limit:       int64(query.Limit),
		Total:       int64(total),
		TotalPages:  int64(totalPages),
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
		Data:        vehicles,
	}

	if err := s.VehicleCache.SetJSON(key, response, 10*time.Minute); err != nil {
		slog.Error("failed to set vehicles in cache", "error", err)
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

func (s *Service) DecreaseStock(ctx context.Context, vehicleID int64, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("invalid quantity")
	}

	exists, err := s.VehicleRepo.ExistByID(int64(vehicleID))
	if err != nil {
		return fmt.Errorf("failed to check existence : %w", err)
	}

	if !exists {
		return fmt.Errorf("vehicle not found")
	}

	if err := s.VehicleRepo.DecreaseStock(ctx, int64(vehicleID), quantity); err != nil {
		return fmt.Errorf("failed to decrease vehicle stock : %w", err)
	}

	return nil
}
