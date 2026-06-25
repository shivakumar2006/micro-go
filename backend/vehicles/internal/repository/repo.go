package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"vehicles/internal/models"
)

type VehicleRepository struct {
	Db *sql.DB
}

func NewVehicleRepo(db *sql.DB) *VehicleRepository {
	return &VehicleRepository{
		Db: db,
	}
}

func (v *VehicleRepository) CreateVehicle(vehicles *models.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := v.Db.ExecContext(ctx, `
		INSERT INTO vehicles(name, model, type, category, created_at)
		VALUES($1, $2, $3, $4, $5)
	`, vehicles.Name, vehicles.Model, vehicles.Type, vehicles.Category, vehicles.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create vehicles : %w", err)
	}

	return nil
}

func (v *VehicleRepository) GetVehicleById(id int64) (*models.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var vehicels models.Vehicle

	err := v.Db.QueryRowContext(ctx, `
		SELECT id, name, model, type, category, created_at
		FROM vehicles
		WHERE id = $1
	`, id).Scan(&vehicels.Id, &vehicels.Name, &vehicels.Model, &vehicels.Type, &vehicels.Category, &vehicels.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle by id : %w", err)
	}

	return &vehicels, nil
}
