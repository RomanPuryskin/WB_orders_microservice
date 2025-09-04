package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/orders_api/internal/models"
)

type OrderRepository interface {
	GetOrderByUID(ctx context.Context, uid uuid.UUID) (*models.Order, error)
	InsertOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	GetAllOrders(ctx context.Context) ([]*models.Order, error)
}
