package cache

import (
	"github.com/gofrs/uuid"
	"github.com/orders_api/internal/models"
)

type Cache interface {
	Get(id uuid.UUID) (*models.Order, bool)
	Set(id uuid.UUID, order *models.Order)
	SetAll([]*models.Order)
}
