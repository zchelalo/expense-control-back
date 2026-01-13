package bootstrap

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zchelalo/expense-control-back/internal/db/connection"
)

func InitDB(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	dsn, err := PostgresDSN(cfg)
	if err != nil {
		return nil, fmt.Errorf("build postgres dsn: %w", err)
	}
	return connection.NewPool(ctx, dsn)
}