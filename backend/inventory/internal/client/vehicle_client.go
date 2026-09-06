package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inventory/internal/metrics"
	"inventory/internal/models"
	"inventory/internal/resilience"
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

func (v *VehicleClient) DecreaseStock(vehicleID int, quantity int) error {
	_, err := resilience.Execute(v.CB, func() (struct{}, error) {
		return resilience.DoWithRetry(v.Retry, func() (struct{}, error) {
			return v.doRequest(vehicleID, quantity)
		})
	})
	return err
}

func (v *VehicleClient) doRequest(vehicleID, quantity int) (struct{}, error) {
	url := fmt.Sprintf("%s/api/v1/vehicles/%d/stock-decrease", v.BaseURL, vehicleID)

	start := time.Now()

	metrics.InventoryVehicleServiceCallTotal.Inc()
	defer func() {
		duration := time.Since(start)
		metrics.InventoryVehicleServiceCallDuration.Observe(duration.Seconds())
	}()

	slog.Info("calling vehicle service", slog.Int("vehicle_id", vehicleID), slog.String("url", url))

	reqBody := models.DecreaseStockRequest{
		Quantity: quantity,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("failed to marshal request", slog.String("error", err.Error()))
		return struct{}{}, fmt.Errorf("failed to marshal request: %v", err)
	}
	slog.Info("request marshaled", slog.Int("vehicle_id", vehicleID))

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		slog.Error("failed to create request", slog.String("error", err.Error()))
		return struct{}{}, fmt.Errorf("failed to create request: %v", err)
	}
	slog.Info("request created", slog.Int("vehicle_id", vehicleID))

	req.Header.Set("Content-Type", "application/json")

	res, err := v.Client.Do(req)
	if err != nil {
		slog.Error("failed to perform request", slog.String("error", err.Error()))
		return struct{}{}, fmt.Errorf("failed to perform request: %v", err)
	}
	slog.Info("request performed", slog.Int("vehicle_id", vehicleID))
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return struct{}{}, &resilience.HTTPStatusError{
			StatusCode: res.StatusCode,
		}
	}
	slog.Info("request performed successfully", slog.Int("vehicle_id", vehicleID))

	return struct{}{}, nil
}
