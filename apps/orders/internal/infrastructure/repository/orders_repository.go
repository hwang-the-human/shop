package repository

import (
	"context"
	"database/sql"
	"errors"
	"shop/apps/orders/domain/entities"
	"shop/shared/repository"
)

type OrdersRepository struct {
	repo repository.Repository
}

func NewOrdersRepository(repo repository.Repository) *OrdersRepository {
	return &OrdersRepository{
		repo: repo,
	}
}

func (r *OrdersRepository) CreateOrder(ctx context.Context, order *entities.Order) error {
	err := r.repo.QueryRow(ctx, "INSERT INTO orders (name, price, quantity, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		order.Name, order.Price, order.Quantity, order.Status, order.CreatedAt, order.UpdatedAt).Scan(&order.ID)
	return err
}

func (r *OrdersRepository) GetOrder(ctx context.Context, id int) (*entities.Order, error) {
	row := r.repo.QueryRow(ctx, "SELECT id, name, price, quantity, status, created_at, updated_at FROM orders WHERE id = $1", id)
	order := &entities.Order{}
	err := row.Scan(&order.ID, &order.Name, &order.Price, &order.Quantity, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return order, err
}

func (r *OrdersRepository) UpdateOrder(ctx context.Context, order *entities.Order) error {
	_, err := r.repo.Exec(ctx, "UPDATE orders SET name = $1, price = $2, quantity = $3, status = $4, updated_at = $5 WHERE id = $6",
		order.Name, order.Price, order.Quantity, order.Status, order.UpdatedAt, order.ID)
	return err
}

func (r *OrdersRepository) DeleteOrder(ctx context.Context, id int) error {
	_, err := r.repo.Exec(ctx, "DELETE FROM orders WHERE id = $1", id)
	return err
}
