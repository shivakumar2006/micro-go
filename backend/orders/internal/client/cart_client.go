package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"orders/internal/models"
	"orders/internal/resilience"
	"time"
)

type CartClient struct {
	BaseURL string `json:"base_url"`
	Client  http.Client
	Retry   *resilience.Retry
}

func NewCartClient(baseURL string) *CartClient {
	return &CartClient{
		BaseURL: baseURL,
		Client:  http.Client{Timeout: 5 * time.Second},
	}
}

func (c *CartClient) GetCart(id int) (*models.CartResponse, error) {
	return resilience.RetryDo(c.Retry, func() (*models.CartResponse, error) {
		return c.doRequest(id)
	})
}

func (c *CartClient) doRequest(id int) (*models.CartResponse, error) {
	url := fmt.Sprintf("%s/api/v1/cart/%d", c.BaseURL, id)

	slog.Info("Requesting cart items", slog.String("url", url))

	start := time.Now()

	res, err := c.Client.Get(url)
	if err != nil {
		slog.Error("failed to call cart service", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to call cart service : %w", err)
	}

	defer res.Body.Close()

	slog.Info("cart service response", slog.Int("status", res.StatusCode), slog.Duration("duration", time.Since(start)))

	if res.StatusCode != http.StatusOK {
		slog.Error("vehicle service returned non-ok status", slog.Int("Status", res.StatusCode))
		return nil, fmt.Errorf("service is not returned the ok status : %w", err)
	}

	var cart models.CartResponse

	if err := json.NewDecoder(res.Body).Decode(&cart); err != nil {
		slog.Error("failed to decode response", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to decode response : %w", err)
	}

	if cart.Status != http.StatusOK {
		slog.Error("cart service returned non-ok status", slog.Int("Status", cart.Status))
		return nil, fmt.Errorf("service is not returned the ok status : %w", err)
	}

	return &cart, nil
}
