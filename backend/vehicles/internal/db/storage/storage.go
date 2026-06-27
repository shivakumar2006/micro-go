package storage

import "vehicles/internal/models"

type Storage interface {
	CreateVehicle(*models.Vehicle) error
	GetVehicleById(id int64) (*models.Vehicle, error)
	GetAllVehicles() ([]models.Vehicle, error)
	UpdateVehicle(*models.Vehicle) error
	DeleteVehicle(id int64) error
	ExistByName(name string) (bool, error)
}
