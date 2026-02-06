package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (*pgxpool.Pool, error) {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		dbName   = os.Getenv("DB_NAME")
	)

	return pgxpool.New(
		context.Background(),
		fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port,
			dbName,
		),
	)
}
