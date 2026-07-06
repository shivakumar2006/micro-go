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

func (c *CartRepository) AddToCart(cart *models.Cart) error {
	ctx, cancel := newContext()
	defer cancel()

	_, err := c.DB.ExecContext(ctx, `
		INSERT INTO cart (user_id, vehicle_id, price, quantity, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6)
	`, cart.UserID, cart.VehicleId, cart.Price, cart.Quantity, cart.CreatedAt, cart.UpdatedAt)

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

func (c *CartRepository) GetUserCart(userId int) ([]models.Cart, error) {
	ctx, cancel := newContext()
	defer cancel()

	var cart []models.Cart

	rows, err := c.DB.QueryContext(ctx, `
		SELECT id, user_id, vehicle_id, price, quantity, created_at, updated_at
		FROM cart 
		WHERE user_id = $1
	`, userId)

	if err != nil {
		return nil, fmt.Errorf("failed to get user cart : %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var cartItems models.Cart
		if err := rows.Scan(
			&cartItems.ID,
			&cartItems.UserID,
			&cartItems.VehicleId,
			&cartItems.Price,
			&cartItems.Quantity,
			&cartItems.CreatedAt,
			&cartItems.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan cart items : %w", err)
		}

		cart = append(cart, cartItems)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get user cart : %w", err)
	}

	return cart, nil
}

func (c *CartRepository) UpdateCartQuantity(cartId int, quantity int) (int, error) {
	ctx, cancel := newContext()
	defer cancel()

	result, err := c.DB.ExecContext(ctx, `
		UPDATE cart
		SET quantity = $1, updated_at = NOW()
		WHERE id = $2
	`, quantity, cartId)
	if err != nil {
		return 0, fmt.Errorf("failed to update cart quantity : %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected : %w", err)
	}

	if rows == 0 {
		return 0, fmt.Errorf("no item found in cart")
	}

	return int(rows), nil
}

func (c *CartRepository) DeleteCartItem(id int) (int, error) {
	ctx, cancel := newContext()
	defer cancel()

	result, err := c.DB.ExecContext(ctx, `
		DELETE from cart 
		WHERE id = $1 
	`, id)

	if err != nil {
		return 0, fmt.Errorf("failed to delete item from cart : %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected : %w", err)
	}

	if rows == 0 {
		return 0, fmt.Errorf("no items in the cart : %w", err)
	}

	return int(rows), nil
}

func (c *CartRepository) DeleteCart(id int) (int, error) {
	ctx, cancel := newContext()
	defer cancel()

	result, err := c.DB.ExecContext(ctx, `
		DELETE FROM cart
		WHERE user_id = $1
	`, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete cart: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return 0, fmt.Errorf("no items in the cart: %w", err)
	}

	return int(rows), nil
}

func (c *CartRepository) ExistByUserAndVehicle(userId, vehicleId int) (bool, error) {
	ctx, cancel := newContext()
	defer cancel()

	var exists bool

	err := c.DB.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 
			FROM cart
			WHERE user_id = $1 AND vehicle_id = $2
		)
	`, userId, vehicleId).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check if item exists in cart : %w", err)
	}

	return exists, nil
}

func (c *CartRepository) CountItems(userId int) (int, error) {
	ctx, cancel := newContext()
	defer cancel()

	var count int

	err := c.DB.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM cart
		WHERE user_id = $1
	`, userId).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count items in cart : %w", err)
	}

	return count, nil
}

func (c *CartRepository) GetCartItem(userId, vehicleId int64) (*models.Cart, error) {
	ctx, cancel := newContext()
	defer cancel()

	var cart models.Cart

	err := c.DB.QueryRowContext(ctx, `
		SELECT id, user_id, vehicle_id, price, quantity, created_at, updated_at
		FROM cart
		WHERE user_id = $1 
		AND vehicle_id = $2
	`, userId, vehicleId).Scan(&cart.ID, &cart.UserID, &cart.VehicleId, &cart.Price, &cart.Quantity, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("cart item not found")
		}
		return nil, fmt.Errorf("failed to get cart items : %w", err)
	}

	return &cart, nil
}

func (c *CartRepository) IncrementQuantity(userId, vehicleId int64, quantity int) (int, error) {
	ctx, cancel := newContext()
	defer cancel()

	result, err := c.DB.ExecContext(ctx, `
		UPDATE cart
		SET quantity = quantity + $1,
			updated_at = NOW()
		WHERE user_id = $2
		AND vehicle_id = $3
	`, quantity, userId, vehicleId)
	if err != nil {
		return 0, fmt.Errorf("failed to increment quantity : %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get row affected: %w", err)
	}

	if rows == 0 {
		return 0, fmt.Errorf("item is not found to increment: %w", err)
	}

	return int(rows), nil
}

func (c *CartRepository) DecrementQuantity(userId, vehicleId int64) (int, error) {
	ctx, cancel := newContext()
	defer cancel()

	var rows int64

	err := c.DB.QueryRowContext(ctx, `
		UPDATE cart
		SET quantity = quantity - 1,
			updated_at = NOW()
		WHERE user_id = $1
		AND vehicle_id = $2
		RETURNING id
	`, userId, vehicleId).Scan(&rows)

	if err != nil {
		return 0, fmt.Errorf("failed to decrement quantity: %w", err)
	}

	if rows == 0 {
		return 0, fmt.Errorf("item is not found in cart to decrement: %w", err)
	}

	return int(rows), nil
}

func (c *CartRepository) GetCartTotal(userId int64) (float64, error) {
	ctx, cancel := newContext()
	defer cancel()

	var total float64

	err := c.DB.QueryRowContext(ctx, `
		SELECT SUM(price * quantity) as total
		FROM cart
		WHERE user_id = $1
	`, userId).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate cart total: %w", err)
	}

	return total, nil
}

func (c *CartRepository) UpdatePrice(cartID int64, price float64) error {
	ctx, cancel := newContext()
	defer cancel()

	_, err := c.DB.ExecContext(ctx, `
		UPDATE cart
		SET price = $1,
			updated_at = NOW()
		WHERE id = $2
	`, cartID, price)
	if err != nil {
		return fmt.Errorf("failed to update price : %w", err)
	}

	return nil
}

func (r *CartRepository) GetCartItemByID(cartID int64) (*models.Cart, error) {
	ctx, cancel := newContext()
	defer cancel()

	var cart models.Cart
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, user_id, vehicle_id, price, quantity, created_at, updated_at
		FROM cart
		WHERE id = $1
	`, cartID).Scan(&cart.ID, &cart.UserID, &cart.VehicleId, &cart.Price, &cart.Quantity, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, fmt.Errorf("item is not found in cart : %w", err)
		}
		return nil, fmt.Errorf("failed to get the cart item by their id : %w", err)
	}

	return &cart, nil
}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
