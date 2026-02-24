package redis

import "example.com/m/internal/infrastructure/env"

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisConfig() Config {
	return Config{
		// Redisのデフォルトポートは6379
		Host:     env.GetString("REDIS_HOST", "localhost"),
		Port:     env.GetString("REDIS_PORT", "6379"),
		Password: env.GetString("REDIS_PASSWORD", ""), // デフォルトは空
		DB:       env.GetInt("REDIS_DB", 0),           // デフォルトのDB番号は0
	}
}

type QueueConfig struct {
	Key string
}
