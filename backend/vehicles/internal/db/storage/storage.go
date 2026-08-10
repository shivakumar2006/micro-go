package storage

import (
	"context"
	"vehicles/internal/models"
)

type Storage interface {
	CreateVehicle(*models.Vehicle) error
	GetVehicleById(id int64) (*models.Vehicle, error)
	GetAllVehicles(models.VehicleQuery) ([]models.Vehicle, int, error)
	UpdateVehicle(*models.Vehicle) error
	DeleteVehicle(id int64) error
	ExistByName(name string) (bool, error)
	ExistsByNameExceptId(name string, id int64) (bool, error)
	ExistByID(id int64) (bool, error)
	DecreaseStock(ctx context.Context, vehicleID int64, quantity int) error
}
