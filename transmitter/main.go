package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"time"
	"transmitter/types"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func main() {
	config := LoadConfig()
	count := flag.Int("count", 10, "Number of loggers to simulate")
	interval := flag.Int("interval", 60, "Interval in seconds between sending data to the server")
	flag.Parse()

	logrus.WithFields(logrus.Fields{
		"Count":    *count,
		"Interval": *interval,
	}).Info("Transmitter started successfully")

	for i := 0; i < *count; i++ {
		go func() {
			seed := rand.IntN(60)
			time.Sleep(time.Duration(seed) * time.Second)

			ticker := time.NewTicker(time.Duration(*interval) * time.Second)
			defer ticker.Stop()

			deviceID := uuid.New().String()

			opts := mqtt.NewClientOptions().
				AddBroker(fmt.Sprintf("tcp://%s:%d", config.MQTT.Address, config.MQTT.Port)).
				SetClientID(fmt.Sprintf("logger-client-%d", i+1)).
				SetCleanSession(true)

			client := mqtt.NewClient(opts)
			for {
				if token := client.Connect(); token.Wait() && token.Error() != nil {
					logrus.Errorf("Failed to connect: %v", token.Error())
					logrus.Info("trying to connect in 10 seconds...")
					time.Sleep(10 * time.Second)
					continue
				}
				break
			}
			defer client.Disconnect(250)

			logrus.Infof("Logger %d successfully connected to EMQX", i+1)

			consumed := 0.0

			for range ticker.C {
				logrus.Info("Send data from logger", i+1)
				consumed += float64(rand.IntN(100)) / 10
				data := types.LoggerData{
					SchemaVersion:   1,
					DeviceID:        deviceID,
					Timestamp:       time.Now().Unix(),
					FirmwareVersion: "1.0.0",
					MeterReading: types.MeterReading{
						Voltage: types.ThreePhase{
							PhaseA: float32(200+rand.IntN(50)) + float32(rand.IntN(10))/10,
							PhaseB: float32(200+rand.IntN(50)) + float32(rand.IntN(10))/10,
							PhaseC: float32(200+rand.IntN(50)) + float32(rand.IntN(10))/10,
						},
						Current: types.ThreePhase{
							PhaseA: float32(rand.IntN(100)),
							PhaseB: float32(rand.IntN(100)),
							PhaseC: float32(rand.IntN(100)),
						},
						Frequency:         50 + float32(rand.IntN(100))/10,
						PowerFactor:       1,
						ActivePowerKW:     500,
						ReactivePowerKVAR: 200,
						ApparentPowerKVA:  800,
						EnergyConsumedKWH: 1234.5 + consumed,
						EnergyExportedKWH: 5,
					},
					Events: []types.Event{
						types.Event{
							Type:      "OVER_VOLTAGE",
							Value:     250.6,
							Timestamp: time.Now().Unix(),
						},
					},
					Location: types.Location{
						Latitude:  23.45663,
						Longitude: 12.436432,
					},
					SignalStrengthDBM:   -67,
					BatteryLevelPercent: 100,
					Status:              "OK",
				}

				payload, err := json.Marshal(data)
				if err != nil {
					log.Printf("error marshalling data: %v", err)
					continue
				}

				token := client.Publish(config.MQTT.Topic, 0, false, payload)
				token.Wait()
			}
		}()
	}

	select {} // block forever
}
