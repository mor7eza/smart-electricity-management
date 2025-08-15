package main

import (
	"billing-service/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type payload struct {
	DeviceID string
	Events   []types.Event
}

func (app *App) Process() {
	var (
		ctx        = context.Background()
		loggerData types.LoggerData
	)
	// get data from RedisService
	for {
		res, err := app.RedisService.Client.BRPop(ctx, 0*time.Second, "telemetry").Result()
		if err != nil {
			logrus.Errorf("error getting data from Redis: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		data := []byte(res[1])
		if err := json.Unmarshal(data, &loggerData); err != nil {
			logrus.Errorf("error unmarshalling json data: %v", err)
		}

		// process consumed energy and save to redis. id, consumed today, totalconsume, timestamp
		loggerDate := time.Unix(loggerData.Timestamp, 0).Format("2006-01-02")
		app.RedisService.Client.HSet(
			ctx,
			fmt.Sprintf("loggers:%s:%s", loggerDate, loggerData.DeviceID),
			map[string]any{
				"timestamp":         loggerData.Timestamp,
				"EnergyConsumedKWH": loggerData.MeterReading.EnergyConsumedKWH,
			},
		)

		// process events and send them to save
		if len(loggerData.Events) > 0 {
			p := payload{
				DeviceID: loggerData.DeviceID,
				Events:   loggerData.Events,
			}

			//send to rabbitmq
			bytes, err := json.Marshal(p)
			if err != nil {
				logrus.Errorf("error marshalling event data: %v", err)
			}
			app.RabbitMQService.PublishMessage(bytes)
		}
	}
}
