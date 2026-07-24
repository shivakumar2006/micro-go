package client

import (
	"cart/internal/models"
	"cart/internal/resilience"
	"encoding/json"
	"fmt"
	"log/slog"
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

	slog.Info("Calling vehicle service", slog.Int("vehicle_id", id), slog.String("url", url))

	start := time.Now()

	res, err := v.Client.Get(url)
	if err != nil {
		slog.Error("Failed to call vehicle service", slog.Int("vehicle_id", id), slog.String("error", err.Error()))
		return nil, err
	}

	defer res.Body.Close()

	slog.Info("Vehicle service response", slog.Int("status_code", res.StatusCode), slog.Duration("duration", time.Since(start)))

	if res.StatusCode != http.StatusOK {

		slog.Error("Vehicle Service returned non-OK status", slog.Int("vehicle_id", id), slog.Int("status_code", res.StatusCode))

		return nil, &resilience.HTTPStatusError{
			StatusCode: res.StatusCode,
		}
	}

	var vehicle models.VehicleResponse

	if err := json.NewDecoder(res.Body).Decode(&vehicle); err != nil {
		slog.Error("failed to decode vehicle service", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to decode response : %w", err)
	}

	return &vehicle, nil
}
