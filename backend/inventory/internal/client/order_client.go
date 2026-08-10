package client

import (
	"encoding/json"
	"fmt"
	"inventory/internal/models"
	"inventory/internal/resilience"
	"log/slog"
	"net/http"
	"time"
)

type OrderClient struct {
	BaseURL string
	Client  *http.Client
	Retry   *resilience.Retry
}

func NewOrderClient(baseUrl string, retry *resilience.Retry) *OrderClient {
	return &OrderClient{
		BaseURL: baseUrl,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Retry: retry,
	}
}

func (c *OrderClient) GetOrders(orderID int) (*models.OrderResponse, error) {
	return resilience.DoWithRetry(c.Retry, func() (*models.OrderResponse, error) {
		return c.doRequest(orderID)
	})
}

func (c *OrderClient) doRequest(orderID int) (*models.OrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/orders/%d", c.BaseURL, orderID)
	if url == "" {
		return nil, fmt.Errorf("order service url not set")
	}
	slog.Info("requesting order items from order service", slog.String("url", url))

	start := time.Now()

	res, err := c.Client.Get(url)
	if err != nil {
		slog.Error("failed to call order service", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get orders from order service : %w", err)
	}

	defer res.Body.Close()

	slog.Info("order service response", slog.Int("status", res.StatusCode), slog.Duration("duration", time.Since(start)))

	if res.StatusCode != http.StatusOK {
		slog.Warn("non 200 status code from order service", slog.String("status", res.Status))
		return nil, &resilience.HTTPStatusError{
			StatusCode: res.StatusCode,
		}
	}

	var OrderResponse models.OrderResponse
	if err := json.NewDecoder(res.Body).Decode(&OrderResponse); err != nil {
		slog.Error("failed to unmarshal order service response")
		return nil, fmt.Errorf("failed to unmarshal order service response : %w", err)
	}

	slog.Info("order service response unmarshalled", slog.Any("order", OrderResponse))
	return &OrderResponse, nil
}
