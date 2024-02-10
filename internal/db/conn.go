package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mauricioabreu/keep/internal/config"
)

func Connect(ctx context.Context, cfg *config.Config) error {
	_, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	return nil
}
