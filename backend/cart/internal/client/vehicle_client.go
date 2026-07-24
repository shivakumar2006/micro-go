package client

import (
	"cart/internal/models"
	"cart/internal/resilience"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type VehicleClient struct {
	BaseURL string
	Client  *http.Client
	Retry   *resilience.Retry
	CB      *resilience.CircuitBreaker
}

func NewVehicleClient(baseURL string, retry *resilience.Retry, cb *resilience.CircuitBreaker) *VehicleClient {
	return &VehicleClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
		Retry:   retry,
		CB:      cb,
	}
}

func (v *VehicleClient) GetVehicle(id int) (*models.VehicleResponse, error) {
	return resilience.Execute(v.CB, func() (*models.VehicleResponse, error) {
		return resilience.RetryDo(v.Retry, func() (*models.VehicleResponse, error) {
			return v.doRequest(id)
		})
	})
}

func (v *VehicleClient) doRequest(id int) (*models.VehicleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/vehicles/%d", v.BaseURL, id)

	res, err := v.Client.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, &resilience.HTTPStatusError{
			StatusCode: res.StatusCode,
		}
	}

	var vehicle models.VehicleResponse

	if err := json.NewDecoder(res.Body).Decode(&vehicle); err != nil {
		return nil, fmt.Errorf("failed to decode response : %w", err)
	}

	return &vehicle, nil
}
