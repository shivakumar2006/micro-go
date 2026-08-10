package service

import (
	"context"
	"fmt"
	"inventory/internal/client"
	"inventory/internal/models"
	"inventory/internal/repository"
)

type InventoryService struct {
	Repo          repository.InventoryRepository
	OrderClient   *client.OrderClient
	VehicleClient *client.VehicleClient
}

func NewInventoryService(repo repository.InventoryRepository, orderClient *client.OrderClient, vehicleClient *client.VehicleClient) *InventoryService {
	return &InventoryService{
		Repo:          repo,
		OrderClient:   orderClient,
		VehicleClient: vehicleClient,
	}
}

func (s *InventoryService) CreateTransaction(ctx context.Context, tx *models.Inventory) error {
	orderId := int(tx.OrderID)

	order, err := s.OrderClient.GetOrders(orderId)
	if err != nil {
		return fmt.Errorf("failed to get order : %w", err)
	}

	for _, item := range order.Items {
		err := s.VehicleClient.DecreaseStock(int(item.VehicleID), item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to decrease stock for vehicle %d : %w", item.VehicleID, err)
		}

		tx := models.Inventory{
			OrderID:   order.ID,
			VehicleID: item.VehicleID,
			Quantity:  item.Quantity,
			Operation: models.OperationDebit,
			Status:    models.StatusCompleted,
		}

		if err := s.Repo.CreateTransaction(ctx, &tx); err != nil {
			return fmt.Errorf("failed to create inventory transaction : %w", err)
		}
	}
	return nil
}
