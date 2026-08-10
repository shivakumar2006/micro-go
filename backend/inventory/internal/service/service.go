package service

import (
	"inventory/internal/client"
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
