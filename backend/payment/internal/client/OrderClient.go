package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"payment/internal/models"
	"time"
)

type OrderClient struct {
	BaseURL string `json:"base_url"`
	Client  *http.Client
}

func NewOrderClient(baseURL string) *OrderClient {
	return &OrderClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (o *OrderClient) GetOrder(orderId int) (*models.OrderResponse, error) {
	return o.doRequest(orderId)
}

func (o *OrderClient) doRequest(orderId int) (*models.OrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/orders/%d", o.BaseURL, orderId)

	slog.Info("Calling vehicle service", slog.Int("order_id : ", orderId), slog.String("url : ", url))

	start := time.Now()

	res, err := o.Client.Get(url)
	if err != nil {
		slog.Error("failed to call order service", slog.Int("order_id : ", orderId), slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to call order service : %w", err)
	}

	defer res.Body.Close()

	slog.Info("order service response", slog.Int("status_code", res.StatusCode), slog.Duration("duration", time.Since(start)))

	if res.StatusCode != http.StatusOK {
		slog.Error("order service returned non-ok status", slog.Int("status_code", res.StatusCode))
		return nil, fmt.Errorf("order service returned non-ok status: %d", res.StatusCode)
	}

	var orderResponse models.OrderResponse

	if err := json.NewDecoder(res.Body).Decode(&orderResponse); err != nil {
		slog.Error("failed to decode order response", slog.Int("order_id : ", orderId), slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to decode order response : %w", err)
	}

	return &orderResponse, nil
}
