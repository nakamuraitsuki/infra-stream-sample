package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"

	"example.com/m/worker"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	
	maxWorkers := getEnvInt("MAX_WORKERS", 5)
	baseURL := getEnv("BASE_URL", "http://frontend")
	worker.StartWorker(ctx, maxWorkers, baseURL)
}
