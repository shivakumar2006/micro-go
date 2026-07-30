package services

import (
	"context"
	"errors"
	"fmt"
	"orders/internal/models"
	"orders/internal/repository"
)

type OrderService struct {
	Repo repository.Repository
}

func NewOrderService(repo repository.Repository) *OrderService {
	return &OrderService{
		Repo: repo,
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

	defer func() {
		tx.Rollback()
	}()

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

	err := s.Repo.UpdateOrderStatus(ctx, int64(orderId), status)
	if err != nil {
		return fmt.Errorf("failed to update order status : %w", err)
	}

	return nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderId int) error {
	if orderId <= 0 {
		return fmt.Errorf("invalid order id")
	}

	if err := s.Repo.UpdateOrderStatus(ctx, int64(orderId), "CANCEL"); err != nil {
		return fmt.Errorf("failed to update the status to cancel of the order id : %w", err)
	}

	return nil
}

func (s *OrderService) MarkOrderpaid(ctx context.Context, orderId int) error {
	if orderId <= 0 {
		return errors.New("invalid order id")
	}

	order, err := s.Repo.GetOrderByID(ctx, int64(orderId))
	if err != nil {
		return fmt.Errorf("failed to get order by id : $%w", err)
	}

	if order.Status == "PENDING" {
		return errors.New("order is not in pending state")
	}

	order.Status = "PAID"

	err = s.Repo.UpdateOrderStatus(ctx, int64(orderId), order.Status)
	if err != nil {
		return fmt.Errorf("failed to update order status")
	}

	return nil
}

// type OrderService interface {

// 	CreateOrder(...)

// 	GetOrderByID(...)

// 	GetOrdersByUserID(...)

// 	UpdateOrderStatus(...)

// 	CancelOrder(...)

// 	MarkOrderPaid(...)

// }
