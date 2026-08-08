package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"payment/internal/models"
	"payment/internal/resilience"
	"time"
)

type OrderClient struct {
	BaseURL        string `json:"base_url"`
	Client         *http.Client
	Retry          *resilience.Retry
	CircuitBreaker *resilience.CircuitBreaker
}

func NewOrderClient(baseURL string, retry *resilience.Retry, cb *resilience.CircuitBreaker) *OrderClient {
	return &OrderClient{
		BaseURL:        baseURL,
		Client:         &http.Client{Timeout: 10 * time.Second},
		Retry:          retry,
		CircuitBreaker: cb,
	}
}

func (o *OrderClient) GetOrder(orderId int) (*models.OrderResponse, error) {
	return resilience.Execute(o.CircuitBreaker, func() (*models.OrderResponse, error) {
		return resilience.DoRetry(o.Retry, func() (*models.OrderResponse, error) {
			return o.doRequest(orderId)
		})
	})
}

func (o *OrderClient) doRequest(orderId int) (*models.OrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/orders/%d", o.BaseURL, orderId)

	slog.Info("Calling order service", slog.Int("order_id", orderId), slog.String("url ", url))

	start := time.Now()

	res, err := o.Client.Get(url)
	if err != nil {
		slog.Error("failed to call order service", slog.Int("order_id", orderId), slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to call order service : %w", err)
	}

	defer res.Body.Close()

	slog.Info("order service response", slog.Int("status_code", res.StatusCode), slog.Duration("duration", time.Since(start)))

	if res.StatusCode != http.StatusOK {
		slog.Error("order service returned non-ok status", slog.Int("status_code", res.StatusCode))
		return nil, &resilience.HTTPStatusError{
			StatusCode: res.StatusCode,
		}
	}

	var orderResponse models.OrderResponse

	if err := json.NewDecoder(res.Body).Decode(&orderResponse); err != nil {
		slog.Error("failed to decode order response", slog.Int("order_id", orderId), slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to decode order response : %w", err)
	}

	return &orderResponse, nil
}
