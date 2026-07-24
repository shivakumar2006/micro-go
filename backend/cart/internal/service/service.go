package service

import (
	"cart/internal/client"
	"cart/internal/db/storage"
	"cart/internal/models"
	cachestorage "cart/internal/redis/cacheStorage"
	"errors"
	"fmt"
	"time"
)

type CartService struct {
	CartRepo      storage.CartStorage
	CartCache     cachestorage.CacheStorage
	VehicleClient *client.VehicleClient
}

func NewCartService(repo storage.CartStorage, cartCache cachestorage.CacheStorage, vehicleClient *client.VehicleClient) *CartService {
	return &CartService{
		CartRepo:      repo,
		CartCache:     cartCache,
		VehicleClient: vehicleClient,
	}
}

func (c *CartService) AddToCart(cart *models.Cart) error {
	defer func() {
		c.CartCache.Delete(fmt.Sprintf("cart:%d", cart.UserID))
		c.CartCache.Delete(fmt.Sprintf("cart:total:%d", cart.UserID))
		c.CartCache.Delete(fmt.Sprintf("cart:count:%d", cart.UserID))
	}()

	// get data from vehicle service
	vehicle, err := c.VehicleClient.GetVehicle(cart.VehicleId)
	if err != nil {
		return err
	}

	if vehicle.Stock <= 0 {
		return fmt.Errorf("vehicle out of stock")
	}

	cart.Price = vehicle.Price

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
	key := fmt.Sprintf("cart:%d", userId)

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

	c.CartCache.Delete(fmt.Sprintf("cart:%d", userId))
	c.CartCache.Delete(fmt.Sprintf("cart:total:%d", userId))
	c.CartCache.Delete(fmt.Sprintf("cart:count:%d", userId))

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

	c.CartCache.Delete(fmt.Sprintf("cart:%d", userID))
	c.CartCache.Delete(fmt.Sprintf("cart:total:%d", userID))
	c.CartCache.Delete(fmt.Sprintf("cart:count:%d", userID))

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

	c.CartCache.Delete(fmt.Sprintf("cart:%d", id))
	c.CartCache.Delete(fmt.Sprintf("cart:total:%d", id))
	c.CartCache.Delete(fmt.Sprintf("cart:count:%d", id))

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
	key := fmt.Sprintf("cart:total:%d", userId)

	var total float64

	if err := c.CartCache.GetJSON(key, &total); err == nil {
		return total, nil
	}

	total, err := c.CartRepo.GetCartTotal(int64(userId))
	if err != nil {
		return 0, nil
	}

	c.CartCache.SetJSON(key, total, 10*time.Minute)

	return total, nil
}

func (c *CartService) CountItems(userId int) (int, error) {
	key := fmt.Sprintf("cart:count:%d", userId)

	var count int

	if err := c.CartCache.GetJSON(key, &count); err == nil {
		return count, nil
	}

	count, err := c.CartRepo.CountItems(userId)
	if err != nil {
		return 0, nil
	}

	c.CartCache.SetJSON(key, count, 10*time.Minute)

	return count, nil
}
