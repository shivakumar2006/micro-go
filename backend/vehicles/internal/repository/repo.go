package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"vehicles/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
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
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return errors.New("vehicle already exist")
			}
		}

		return fmt.Errorf("failed to create vehicles : %w", err)
	}

	return nil
}

func (v *VehicleRepository) GetVehicleById(id int64) (*models.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var vehicle models.Vehicle

	err := v.Db.QueryRowContext(ctx, `
		SELECT id, name, model, type, category, created_at
		FROM vehicles
		WHERE id = $1
	`, id).Scan(&vehicle.Id, &vehicle.Name, &vehicle.Model, &vehicle.Type, &vehicle.Category, &vehicle.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("vehicle not found")
		}
		return nil, fmt.Errorf("failed to get vehicle by id : %w", err)
	}

	return &vehicle, nil
}

func (v *VehicleRepository) GetAllVehicles(page, limit int) ([]models.Vehicle, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var total int

	err := v.Db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM vehicles
	`).Scan(&total)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total vehicles : %w", err)
	}

	offset := (page - 1) * limit

	rows, err := v.Db.QueryContext(ctx, `
		SELECT id, name, model, type, category, created_at
		FROM vehicles 
		ORDER BY id 
		LIMIT $1
		OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all vehicles : %w", err)
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var vehicle models.Vehicle
		err := rows.Scan(
			&vehicle.Id,
			&vehicle.Name,
			&vehicle.Model,
			&vehicle.Type,
			&vehicle.Category,
			&vehicle.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed during rows iteration: %w", err)
	}

	return vehicles, total, nil
}

func (v *VehicleRepository) UpdateVehicle(vehicles *models.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := v.Db.ExecContext(ctx, `
		UPDATE vehicles
		SET name = $1, model = $2, type = $3, category = $4
		WHERE id = $5
	`, vehicles.Name, vehicles.Model, vehicles.Type, vehicles.Category, vehicles.Id)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return errors.New("vehicle already exist")
			}
		}
		return fmt.Errorf("failed to update vehicle : %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected : %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("vehicle not found")
	}

	return nil
}

func (v *VehicleRepository) DeleteVehicle(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := v.Db.ExecContext(ctx, `
		DELETE FROM vehicles
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("failed to delete vehicle : %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected : %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("vehicle not found")
	}

	return nil
}

func (v *VehicleRepository) ExistByName(name string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exist bool

	err := v.Db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM vehicles WHERE name = $1)
	`, name).Scan(&exist)

	if err != nil {
		return false, fmt.Errorf("failed to check existence : %w", err)
	}

	return exist, nil
}

func (v *VehicleRepository) ExistsByNameExceptId(name string, id int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exist bool

	err := v.Db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM vehicles WHERE name = $1 AND id != $2)
	`, name, id).Scan(&exist)

	if err != nil {
		return false, fmt.Errorf("failed to check existence : %w", err)
	}

	return exist, nil
}
