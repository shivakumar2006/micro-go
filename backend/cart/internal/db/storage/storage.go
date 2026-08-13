package storage

import "cart/internal/models"

type CartStorage interface {
	AddToCart(cart *models.Cart) error
	GetUserCart(userId int) ([]models.Cart, error)
	UpdateCartQuantity(userID, cartID, quantity int) (int, error)
	DeleteCartItem(userId, cartID int) (int, error)
	DeleteCart(userID int) (int, error)
	ExistByUserAndVehicle(userId, vehicleId int) (bool, error)
	CountItems(userId int) (int, error)
	GetCartItem(userId, vehicleId int64) (*models.Cart, error)
	GetCartItemByID(cartID int64) (*models.Cart, error)
	IncrementQuantity(userId, vehicleId int64, quantity int) (int, error)
	DecrementQuantity(userId, vehicleId int64) (int, error)
	GetCartTotal(userId int64) (float64, error)
	UpdatePrice(cartID int64, price float64) error
}
