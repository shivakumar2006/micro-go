package client

import (
	"fmt"
	"inventory/internal/models"
	"log/slog"
	"net/http"
	"time"
)

type OrderClient struct {
	BaseURL string
	Client  *http.Client
}

func NewOrderClient(baseUrl string) *OrderClient {
	return &OrderClient{
		BaseURL: baseUrl,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *OrderClient) GetOrders(orderID int) (any, error) {
	return c.doRequest(orderID)
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

}
