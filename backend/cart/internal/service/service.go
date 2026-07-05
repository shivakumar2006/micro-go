package service

import (
	"cart/internal/db/storage"
	"cart/internal/models"
	"errors"
)

type CartService struct {
	CartRepo storage.CartStorage
}

func NewCartService(repo storage.CartStorage) *CartService {
	return &CartService{
		CartRepo: repo,
	}
}

func (c *CartService) AddToCart(cart *models.Cart) error {
	if err := c.CheckDuplicate(cart); err != nil {
		return err
	}

	if err := c.CartRepo.AddToCart(cart); err != nil {
		return err
	}
	return nil
}

func (c *CartService) GetUserCart(userId int) ([]models.Cart, error) {
	cart, err := c.CartRepo.GetUserCart(userId)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (c *CartService) CheckDuplicate(cart *models.Cart) error {
	exist, err := c.CartRepo.ExistByUserAndVehicle(cart.UserID, cart.VehicleId)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("cart item already exist")
	}

	return nil
}
