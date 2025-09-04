package cache

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/orders_api/internal/models"
)

type OrderCacher struct {
	store map[uuid.UUID]*models.Order
	mu    sync.RWMutex
}

func NewOrderCacher() *OrderCacher {
	return &OrderCacher{
		store: make(map[uuid.UUID]*models.Order),
	}
}

func (c *OrderCacher) Get(id uuid.UUID) (*models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if order, ok := c.store[id]; ok {
		return order, true
	}
	return nil, false
}

func (c *OrderCacher) Set(id uuid.UUID, order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[id] = order
}

func (c *OrderCacher) SetAll(orders []*models.Order) {
	for _, order := range orders {
		c.Set(order.OrderUID, order)
	}
}
