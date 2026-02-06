package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	var (
		host     = os.Getenv("RDB_HOST")
		port     = os.Getenv("RDB_PORT")
		user     = os.Getenv("RDB_USERNAME")
		password = os.Getenv("RDB_PASSWORD")
		dbName   = os.Getenv("RDB_NAME")
	)

	db, err := strconv.Atoi(dbName)
	if err != nil {
		return nil, fmt.Errorf("invalid RDB_NAME: %w", err)
	}

	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: user,
		Password: password,
		DB:       db,
	}), nil
}
