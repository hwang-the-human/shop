package application

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"shop/apps/orders/domain/entities"
	"shop/apps/orders/infrastructure/repository"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidOrder  = errors.New("invalid order")
)

type OrdersUsecase struct {
	repo *repository.OrdersRepository
}

func NewOrdersService(repo *repository.OrdersRepository) *OrdersUsecase {
	return &OrdersUsecase{
		repo: repo,
	}
}

func (s *OrdersUsecase) CreateOrder(ctx context.Context, name string, price float64, quantity int) (*entities.Order, error) {
	if name == "" || price <= 0 || quantity <= 0 {
		return nil, ErrInvalidOrder
	}

	order := &entities.Order{
		Name:      name,
		Price:     price,
		Quantity:  quantity,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrdersUsecase) GetOrder(ctx context.Context, id int) (*entities.Order, error) {
	order, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}

func (s *OrdersUsecase) UpdateOrder(ctx context.Context, order *entities.Order) error {
	if order == nil || order.ID == 0 {
		return ErrInvalidOrder
	}

	order.UpdatedAt = time.Now()

	return s.repo.UpdateOrder(ctx, order)
}

func (s *OrdersUsecase) DeleteOrder(ctx context.Context, id int) error {
	if _, err := s.GetOrder(ctx, id); err != nil {
		return err
	}

	return s.repo.DeleteOrder(ctx, id)
}
