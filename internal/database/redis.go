package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, host, port string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})

	for i := 0; i < 30; i++ {
		if err := client.Ping(ctx).Err(); err == nil {
			log.Println("✅ Redis connected")
			return client, nil
		}
		log.Printf(" Waiting for Redis... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("cannot connect to Redis after retries")
}
