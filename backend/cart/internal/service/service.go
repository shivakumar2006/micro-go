package service

import (
	"cart/internal/db/storage"
	"cart/internal/models"
	cachestorage "cart/internal/redis/cacheStorage"
	"errors"
	"fmt"
	"time"
)

type CartService struct {
	CartRepo  storage.CartStorage
	CartCache cachestorage.CacheStorage
}

func NewCartService(repo storage.CartStorage, cartCache cachestorage.CacheStorage) *CartService {
	return &CartService{
		CartRepo:  repo,
		CartCache: cartCache,
	}
}

func (c *CartService) AddToCart(cart *models.Cart) error {
	exist, err := c.CartRepo.ExistByUserAndVehicle(cart.UserID, cart.VehicleId)
	if err != nil {
		return err
	}

	if exist {
		_, err := c.CartRepo.IncrementQuantity(int64(cart.UserID), int64(cart.VehicleId), cart.Quantity)
		if err != nil {
			return err
		}

		return nil
	}

	if err := c.CartRepo.AddToCart(cart); err != nil {
		return err
	}
	return nil
}

func (c *CartService) GetUserCart(userId int) ([]models.Cart, error) {
	key := fmt.Sprintf("cart: %d", userId)

	var cart []models.Cart

	if err := c.CartCache.GetJSON(key, &cart); err == nil {
		return cart, nil
	}

	cart, err := c.CartRepo.GetUserCart(userId)
	if err != nil {
		return nil, err
	}

	c.CartCache.SetJSON(key, cart, 10*time.Minute)

	return cart, nil
}

func (c *CartService) UpdateCartQuantity(userId, cartId int, quantity int) (int, error) {
	if quantity <= 0 {
		return 0, errors.New("quantity must be greater than 0")
	}

	updateCart, err := c.CartRepo.UpdateCartQuantity(cartId, quantity)
	if err != nil {
		return 0, err
	}

	return updateCart, nil
}

func (c *CartService) DeleteCartItem(userID, cartID int) (int, error) {
	item, err := c.GetCartItemByID(int64(cartID))
	if err != nil {
		return 0, err
	}

	if item.UserID != userID {
		return 0, errors.New("access denied")
	}

	if item.Quantity <= 0 {
		return 0, errors.New("item is not found in cart")
	}

	if item.Quantity > 1 {
		_, err := c.CartRepo.DecrementQuantity(int64(item.UserID), int64(item.VehicleId))
		if err != nil {
			return 0, err
		}

		return item.Quantity - 1, nil
	}

	deleteItem, err := c.CartRepo.DeleteCartItem(item.UserID, item.ID)
	if err != nil {
		return 0, err
	}

	return deleteItem, nil
}

func (c *CartService) GetCartItemByID(cartID int64) (*models.Cart, error) {
	return c.CartRepo.GetCartItemByID(cartID)
}

func (c *CartService) DeleteCart(id int) (int, error) {
	deleteCart, err := c.CartRepo.DeleteCart(id)
	if err != nil {
		return 0, err
	}
	return deleteCart, nil
}

func (c *CartService) GetCartItem(userId, vehicleId int64) (*models.Cart, error) {
	items, err := c.CartRepo.GetCartItem(userId, vehicleId)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *CartService) GetCartTotal(userId int) (float64, error) {
	return c.CartRepo.GetCartTotal(int64(userId))
}

func (c *CartService) CountItems(userId int) (int, error) {
	return c.CartRepo.CountItems(userId)
}
