package client

import (
	"fmt"
	"net/http"
	"payment/internal/models"
	"payment/internal/resilience"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v78"
	checkoutsession "github.com/stripe/stripe-go/v78/checkout/session"
)

type StripeClient struct {
	BaseURL    string
	SecretKey  string
	SuccessURL string
	CancelURL  string
	Client     *http.Client
	Retry      *resilience.Retry
	CB         *resilience.CircuitBreaker
}

func NewStripeClient(baseURL, secretKey, successURL, cancelURL string, retry *resilience.Retry, cb *resilience.CircuitBreaker) *StripeClient {
	return &StripeClient{
		BaseURL:    baseURL,
		SecretKey:  secretKey,
		SuccessURL: successURL,
		CancelURL:  cancelURL,
		Client:     &http.Client{Timeout: 10 * time.Second},
		Retry:      retry,
		CB:         cb,
	}
}

func (s *StripeClient) CreateCheckoutSession(order *models.OrderResponse) (*models.StripeCheckoutResponse, error) {
	return resilience.Execute(s.CB, func() (*models.StripeCheckoutResponse, error) {
		return resilience.DoRetry(s.Retry, func() (*models.StripeCheckoutResponse, error) {
			stripe.Key = s.SecretKey

			params := &stripe.CheckoutSessionParams{
				PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
				Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
				SuccessURL:         stripe.String(s.SuccessURL),
				CancelURL:          stripe.String(s.CancelURL),
				LineItems: []*stripe.CheckoutSessionLineItemParams{
					{
						PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
							Currency: stripe.String("inr"),

							ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
								Name: stripe.String(fmt.Sprintf("Order #%d", order.ID)),
							},

							UnitAmount: stripe.Int64(int64(order.TotalAmount * 100)),
						},

						Quantity: stripe.Int64(1),
					},
				},

				Metadata: map[string]string{
					"order_id": strconv.Itoa(order.ID),
					"user_id":  strconv.Itoa(order.UserID),
				},
			}
			session, err := checkoutsession.New(params)
			if err != nil {
				return nil, fmt.Errorf("failed to create checkout session : %w", err)
			}

			return &models.StripeCheckoutResponse{
				ID:  session.ID,
				URL: session.URL,
			}, nil
		})
	})
}
