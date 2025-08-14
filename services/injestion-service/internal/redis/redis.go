package redis_db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func NewClient(url string) *redis.Client {
	ctx := context.Background()
	for {
		client := redis.NewClient(&redis.Options{
			Addr: url,
			DB:   0,
		})

		err := client.Ping(ctx).Err()
		if err == nil {
			logrus.Info("Connected to Redis")
			return client
		}

		logrus.Errorf("failed to connect to Redis: %v", err)
		logrus.Info("Retrying to connect in 10s...")
		time.Sleep(10 * time.Second)
	}
}

func Publish(ctx context.Context, c *redis.Client, streamName string, payload []byte) {
	c.Set(ctx, "key string", "value interface{}", 0)
	c.LPush(ctx, "telemetry", payload)
}
