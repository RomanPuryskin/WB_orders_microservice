package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/orders_api/api/handlers"
)

func InitRoutesForOrders(app *fiber.App, handler *handlers.OrderHandler) {
	api := app.Group("/")

	api.Get("/orders/:order_uid", handler.GetOrderByUID)
}

func InitRouteForSwagger(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}
