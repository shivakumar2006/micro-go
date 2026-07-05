package repository

import (
	"cart/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type CartRepository struct {
	DB *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{DB: db}
}

// POST    /cart              // Add to cart
func (c *CartRepository) AddToCart(cart *models.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.DB.ExecContext(ctx, `
		INSERT INTO cart (id, user_id, vehicle_id, price, quantity, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, %6, %7)
	`, cart.ID, cart.UserID, cart.VehicleId, cart.Price, cart.Quantity, cart.CreatedAt, cart.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("item already exists in cart")
			}
		}
		return fmt.Errorf("failed to add item in cart : %w", err)
	}

	return nil
}

// GET     /cart              // User ki cart
func (c *CartRepository) GetUserCart() ([]models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

}

// PUT     /cart/{id}         // Quantity update

// DELETE  /cart/{id}         // Ek item remove

// DELETE  /cart              // Puri cart clear
