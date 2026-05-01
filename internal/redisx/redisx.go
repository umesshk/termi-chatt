package redisx

import (
	"context"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Client struct {
	Rdb *redis.Client
}

func New(addr, password, dbStr string) (*Client, error) {
	if addr == "" {
		return nil, nil
	}

	db := 0
	if dbStr != "" {
		if parsed, err := strconv.Atoi(dbStr); err == nil {
			db = parsed
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return &Client{Rdb: rdb}, nil
}

