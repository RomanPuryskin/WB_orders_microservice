package main

// @title WB_order API
// @version 1.0
import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/orders_api/api/app"
	_ "github.com/orders_api/docs"
	"github.com/orders_api/internal/config"
	"github.com/orders_api/internal/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg, err := config.MustLoad()
	if err != nil {
		slog.Error("Failed config",
			"error", err)
		os.Exit(1)
	}

	// подключим логгер
	logger.InitLogger(&cfg.Logger)

	app := app.InitNewFiberApp(cfg, context.Background())

	app.Start(ctx)

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = app.Stop(ctx)
	if err != nil {
		slog.Error("Failed shutdown", "error", err)
		os.Exit(1)
	}
	slog.Info("Gracefully stopped")
}
