package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/orders_api/api/handlers"
	"github.com/orders_api/api/routes"
	"github.com/orders_api/internal/config"
	"github.com/orders_api/internal/database/cache"
	"github.com/orders_api/internal/database/postgres"
	"github.com/orders_api/internal/kafka"
	"github.com/orders_api/internal/repository"
	"github.com/orders_api/internal/service"
	segmentio "github.com/segmentio/kafka-go"
)

type App struct {
	FiberApp    *fiber.App
	Srvc        service.ServiceOrder
	Db          *pgx.Conn
	Cfg         *config.Config
	KafkaReader *segmentio.Reader
	Consumer    *kafka.KafkaConsumer
}

func InitNewFiberApp(cfg *config.Config, ctx context.Context) *App {

	// подключаемся к БД
	db, err := postgres.NewPostgresDB(ctx, &cfg.Postgres)
	if err != nil {
		slog.Error("Failed connect to postgres DB",
			"error", err)
		os.Exit(1)
	}
	slog.Info("Successfully connected to postgres DB")

	// запускаем миграции
	err = postgres.RunMigrations(&cfg.Postgres)
	if err != nil {
		slog.Error("Failed run migrations",
			"error", err)
		os.Exit(1)
	}
	slog.Info("Successfully ran migratons")

	// создаем репозиторий и кэш для сервиса
	repOrder := repository.NewOrderPostgresRepository(db)
	cacheOrder := cache.NewOrderCacher()

	// создаем сервис обработки заказов
	serviceOrder := service.NewServiceOrder(repOrder, cacheOrder, ctx)

	// при старте сервиса загрузим все актуальные данные из БД в кэш
	err = serviceOrder.Recover()
	if err != nil {
		slog.Info("Started with empty cache")
	} else {
		slog.Info("Successfully loaded data to cache")
	}

	// подключим kafka Reader
	kafkaReader, err := kafka.NewReader(&cfg.Kafka)
	if err != nil {
		slog.Error("Failed run kafka",
			"error", err)
		os.Exit(1)
	}
	slog.Info("Successfully ran kafka", "brokers", cfg.Kafka.Address, "topic", cfg.Kafka.Topic, "group", cfg.Kafka.Group)

	// подключим Consumer
	consumer := kafka.NewKafkaConsumer(kafkaReader, &cfg.Kafka, serviceOrder)

	// создаем новый FiberApp
	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	// подключим хэндлер заказов
	orderHandler := handlers.NewOrderHandler(serviceOrder)

	// подключаем роуты
	routes.InitRoutesForOrders(app, orderHandler)
	routes.InitRouteForSwagger(app)

	app.Static("/", "assets")

	return &App{
		FiberApp:    app,
		Db:          db,
		Srvc:        serviceOrder,
		Cfg:         cfg,
		KafkaReader: kafkaReader,
		Consumer:    consumer,
	}
}

func (a *App) Start(ctx context.Context) {
	slog.Info("App staring", "port", a.Cfg.ServerPort)

	go func() {
		err := a.FiberApp.Listen(fmt.Sprintf(":%s", a.Cfg.ServerPort))
		if err != nil {
			slog.Error("Failed to start app",
				"error", err)
			os.Exit(1)
		}
	}()

	go a.Consumer.ReadMessages(ctx)
	slog.Info("Consumer started")

}

func (a *App) Stop(ctx context.Context) error {
	slog.Info("[!] Shutting down...")

	var stopErr = errors.New("")

	// закрываем соединение БД
	if err := postgres.ClosePostgresDB(ctx, a.Db); err != nil {
		errors.Join(stopErr, err)
	}

	// закрываем kafky
	if err := a.KafkaReader.Close(); err != nil {
		errors.Join(stopErr, err)
	}

	// закрываем соединение с сервером
	if err := a.FiberApp.ShutdownWithContext(ctx); err != nil {
		errors.Join(stopErr, err)
	}

	if stopErr.Error() != "" {
		return stopErr
	}
	return nil
}
