package services

import (
	"context"
	"errors"
	"fmt"
	"orders/internal/client"
	"orders/internal/models"
	"orders/internal/repository"
)

type OrderService struct {
	Repo       repository.Repository
	CartClient *client.CartClient
}

func NewOrderService(repo repository.Repository, cartClient *client.CartClient) *OrderService {
	return &OrderService{
		Repo:       repo,
		CartClient: cartClient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *models.CreateOrderRequest) (*models.Order, error) {
	if req.UserID == 0 {
		return nil, errors.New("invalid user id")
	}

	if len(req.Items) == 0 {
		return nil, errors.New("order items is required")
	}

	for _, item := range req.Items {
		if item.Quantity <= 0 || item.Price <= 0 {
			return nil, errors.New("invalid quantity or price")
		}
	}

	// calc total amount
	var total float64
	for _, item := range req.Items {
		total += float64(item.Quantity) * item.Price
	}

	order := &models.Order{
		UserID:      int(req.UserID),
		Status:      "PENDING",
		TotalAmount: total,
	}

	tx, err := s.Repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transactions : %w", err)
	}

	defer func() {
		tx.Rollback()
	}()

	err = s.Repo.CreateOrder(ctx, tx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order : %w", err)
	}

	var orderItems []models.OrderItem

	for _, item := range req.Items {
		orderItems = append(orderItems, models.OrderItem{
			OrderID:   order.ID,
			VehicleID: item.VehicleID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	err = s.Repo.CreateOrderItems(ctx, tx, orderItems)
	if err != nil {
		return nil, fmt.Errorf("failed to create order items : %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("faild to commit or save : %w", err)
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id int) (*models.Order, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	order, err := s.Repo.GetOrderByID(ctx, int64(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id : %w", err)
	}

	return order, nil
}

func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error) {
	if userID <= 0 {
		return nil, errors.New("invalid id")
	}

	order, err := s.Repo.GetOrdersByUserID(ctx, int64(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get order by user id : %w", err)
	}

	return order, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderId int, status string) error {
	if orderId <= 0 {
		return errors.New("invalid order id")
	}

	switch status {
	case models.OrderStatusPending,
		models.OrderStatusPaid,
		models.OrderStatusCancelled,
		models.OrderStatusDelivered:

	default:
		return errors.New("invalid order status")
	}

	err := s.Repo.UpdateOrderStatus(ctx, int64(orderId), status)
	if err != nil {
		return fmt.Errorf("failed to update order status : %w", err)
	}

	return nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID int64) error {
	if orderID <= 0 {
		return errors.New("invalid order id")
	}

	order, err := s.Repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	switch order.Status {
	case models.OrderStatusCancelled:
		return errors.New("order is already cancelled")

	case models.OrderStatusDelivered:
		return errors.New("delivered order cannot be cancelled")
	}

	if err := s.Repo.UpdateOrderStatus(ctx, orderID, models.OrderStatusCancelled); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

func (s *OrderService) MarkOrderPaid(ctx context.Context, orderID int64) error {
	if orderID <= 0 {
		return errors.New("invalid order id")
	}

	order, err := s.Repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status != models.OrderStatusPending {
		return fmt.Errorf(
			"only pending orders can be marked as paid (current status: %s)",
			order.Status,
		)
	}

	if err := s.Repo.UpdateOrderStatus(ctx, orderID, models.OrderStatusPaid); err != nil {
		return fmt.Errorf("failed to mark order as paid: %w", err)
	}

	return nil
}
