package service

import (
	"inventory/internal/repository"
	"payment/internal/client"
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
