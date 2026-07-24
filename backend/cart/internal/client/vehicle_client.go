package client

import (
	"cart/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type VehicleClient struct {
	BaseURL string
	Client  *http.Client
}

func NewVehicleClient(baseURL string) *VehicleClient {
	return &VehicleClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (v *VehicleClient) GetVehicle(id int) (*models.VehicleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/vehicles/%d", v.BaseURL, id)

	res, err := v.Client.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vehicle not found : %w", err)
	}

	var vehicle models.VehicleResponse

	if err := json.NewDecoder(res.Body).Decode(&vehicle); err != nil {
		return nil, fmt.Errorf("failed to decode response : %w", err)
	}

	return &vehicle, nil
}
