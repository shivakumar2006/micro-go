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

func (s *InventoryService) CreateTransaction(ctx context.Context, orderID int64) error {
	order, err := s.OrderClient.GetOrders(int(orderID))
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

func (s *InventoryService) GetTransactionByID(ctx context.Context, id int64) (*models.Inventory, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}

	result, err := s.Repo.GetTransactionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by id :  %w", err)
	}

	return result, nil
}

func (s *InventoryService) GetTransactionByOrderID(ctx context.Context, orderId int64) ([]models.Inventory, error) {
	if orderId <= 0 {
		return nil, fmt.Errorf("invalid orderid")
	}

	result, err := s.Repo.GetTransactionByOrderID(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by order id : %w", err)
	}

	return result, nil
}

func (s *InventoryService) UpdateTransactionStatus(ctx context.Context, transactionID int64, status string) error {
	if transactionID <= 0 {
		return fmt.Errorf("invalid transactionID")
	}

	exists, err := s.GetTransactionByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("transaction is not exist : %w", err)
	}

	if exists.Status == models.StatusFailed {
		return fmt.Errorf("transaction already failed")
	}

	if exists.Status == models.StatusCompleted {
		return fmt.Errorf("trnasaction already completed")
	}

	if err := s.Repo.UpdateTransactionStatus(ctx, transactionID, status); err != nil {
		return fmt.Errorf("failed to udpate transaction status : %w", err)
	}

	return nil
}

func (s *InventoryService) GetTransactions(ctx context.Context) ([]models.Inventory, error) {
	transactions, err := s.Repo.GetTransactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions : %w", err)
	}

	return transactions, nil
}
