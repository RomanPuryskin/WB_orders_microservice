package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func NewPostgresDB(ctx context.Context, cfg *PostgresConfig) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("[NewPostgresDB|connect] , %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("[NewPostgresDB|ping] , %w", err)
	}

	return conn, nil
}

func ClosePostgresDB(ctx context.Context, conn *pgx.Conn) error {
	err := conn.Close(ctx)
	if err != nil {
		return err
	}
	return nil
}
