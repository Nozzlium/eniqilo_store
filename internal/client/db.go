package client

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/config"
)

func InitDB(cfg config.DBConfig) (*pgx.Conn, error) {
	dbURI := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBParams,
	)

	conn, err := pgx.Connect(
		context.Background(),
		dbURI,
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
