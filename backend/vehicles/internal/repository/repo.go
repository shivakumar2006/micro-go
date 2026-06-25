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

func (v *VehicleRepository) GetAllVehicles() ([]models.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := v.Db.QueryRowContext(ctx, `
		SELECT id, name, model, type, category, created_at
		FROM vehicles 
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to get all vehicles : %w", err)
	}

	return []models.Vehicle{}, nil
}

func (v *VehicleRepository) UpdateVehicle(vehicles *models.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := v.Db.ExecContext(ctx, `
		UPDATE vehicles
		SET name = $1, model = $2, type = $3, category = $4
		WHERE id = $5
	`, vehicles.Name, vehicles.Model, vehicles.Type, vehicles.Category, vehicles.Id)

	if err != nil {
		return fmt.Errorf("failed to update vehicle : %w", err)
	}

	return nil
}

func (v *VehicleRepository) DeleteVehicle(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := v.Db.ExecContext(ctx, `
		DELETE FROM vehicles
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("failed to delete vehicle : %w", err)
	}

	return nil
}
