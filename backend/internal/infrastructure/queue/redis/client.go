package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	Universal redis.UniversalClient
}

func NewClient(cfg Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{Universal: rdb}, nil
}

func (c *Client) Close() error {
	return c.Universal.Close()
}
