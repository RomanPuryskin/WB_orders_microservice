package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/orders_api/internal/database/cache"
	"github.com/orders_api/internal/models"
	"github.com/orders_api/internal/repository"
	"github.com/orders_api/internal/utils"
)

var (
	ErrInvalidUUID  = errors.New("invalid uuid order's fromat")
	ErrValidateJSON = errors.New("invalid values in JSON request")
)

type ServiceOrder interface {
	GetOrderByUID(id string) (*models.Order, error)
	SetOrder(order *models.Order) (*models.Order, error)
	Recover() error
}

type serviceOrder struct {
	Repo  repository.OrderRepository
	Cache cache.Cache
	ctx   context.Context
}

func NewServiceOrder(r repository.OrderRepository, c cache.Cache, ct context.Context) *serviceOrder {
	return &serviceOrder{
		Repo:  r,
		Cache: c,
		ctx:   ct,
	}
}

func (s *serviceOrder) GetOrderByUID(id string) (*models.Order, error) {

	// провалидировать uid
	order_uuid, err := utils.ValidateUUID(id)
	if err != nil {
		return nil, fmt.Errorf("[GetOrderByUID|validate]: %w", ErrInvalidUUID)
	}

	// сначала обращаемся к кэшу
	respOrder, exist := s.Cache.Get(order_uuid)
	if exist {
		slog.Info("got order from cache", "order_uuid", order_uuid)
		return respOrder, nil
	}

	// запрос к БД, если в кэше нет
	respOrder, err = s.Repo.GetOrderByUID(s.ctx, order_uuid)
	if err != nil {
		return nil, err
	}

	// запишем этот заказ в кэш
	s.Cache.Set(order_uuid, respOrder)
	slog.Info("set order into cache", "order_uuid", order_uuid)
	return respOrder, nil

}

func (s *serviceOrder) SetOrder(order *models.Order) (*models.Order, error) {

	// провалидируем структуру на ограничения полей
	err := utils.VaildateStructs(order)
	if err != nil {
		return nil, fmt.Errorf("[SetOrder|validate JSON]: %w", ErrValidateJSON)
	}

	// запрос к БД
	newOrder, err := s.Repo.InsertOrder(s.ctx, order)
	if err != nil {
		return nil, err
	}
	return newOrder, nil
}

func (s *serviceOrder) Recover() error {
	orders, err := s.Repo.GetAllOrders(s.ctx)
	if err != nil {
		return fmt.Errorf("[Recoder| getAllOrders]: %w", err)
	}

	s.Cache.SetAll(orders)

	return nil
}
