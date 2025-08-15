package redis_db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RedisService struct {
	Client *redis.Client
}

func NewService() *RedisService {
	url := viper.GetString("REDIS_URL")
	ctx := context.Background()

	for {
		client := redis.NewClient(&redis.Options{
			Addr: url,
			DB:   0,
		})

		err := client.Ping(ctx).Err()
		if err == nil {
			logrus.Info("Connected to Redis")
			return &RedisService{
				Client: client,
			}
		}

		logrus.Errorf("failed to connect to Redis: %v", err)
		logrus.Info("Retrying to connect in 10s...")
		time.Sleep(10 * time.Second)
	}
}

func (rs *RedisService) Publish(ctx context.Context, streamName string, payload []byte) {
	err := rs.Client.LPush(ctx, "telemetry", payload).Err()
	if err != nil {
		logrus.Errorf("error pushing to redis: %v", err)
	}
}
