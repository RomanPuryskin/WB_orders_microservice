package handlers

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/orders_api/api/errs"
	"github.com/orders_api/internal/repository"
	"github.com/orders_api/internal/service"
)

type OrderHandler struct {
	service service.ServiceOrder
}

func NewOrderHandler(s service.ServiceOrder) *OrderHandler {
	return &OrderHandler{
		service: s,
	}
}

// GetOrderByUID godoc
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя
// @Tags orders
// @Accept json
// @Produce json
// @Param order_uid path string true "Order UUID" Format(uuid)
// @Success 200 {object} models.Order
// @Failure 400 {object} errs.ErrorResponse
// @Failure 404 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByUID(c *fiber.Ctx) error {
	order_id := c.Params("order_uid")

	respOrder, err := h.service.GetOrderByUID(order_id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUUID):
			slog.Error("invalid order_uuid format", "order_uuid", order_id)
			return c.Status(errs.ErrInvalidUUID.Code).JSON(errs.ErrInvalidUUID)

		case errors.Is(err, repository.ErrOrderNotFoundByUUID):
			slog.Error("order not found with order_uuid", "order_uuid", order_id)
			return c.Status(errs.ErrOrderNotFound.Code).JSON(errs.ErrOrderNotFound)

		default:
			slog.Error("error while finding order",
				"order_uuid", order_id,
				"error", err)
			return c.Status(errs.ErrInternalServer.Code).JSON(errs.ErrInternalServer)
		}
	}

	slog.Info("success found order with order_uuid",
		"order_uuid", order_id)
	return c.Status(fiber.StatusOK).JSON(respOrder)

}
