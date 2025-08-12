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
)

func main() {
	config := LoadConfig()
	count := flag.Int("count", 10, "Number of loggers to simulate")
	interval := flag.Int("interval", 60, "Interval in seconds between sending data to the server")
	flag.Parse()

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

			for range ticker.C {
				fmt.Println("Meter", i+1, "-> Sent")
				data := types.LoggerData{
					SchemaVersion:   1,
					DeviceID:        deviceID,
					Timestamp:       time.Now().Unix(),
					FirmwareVersion: "1.0.0",
					MeterReading: types.MeterReading{
						Voltage: types.ThreePhase{
							PhaseA: 220.5,
							PhaseB: 220.5,
							PhaseC: 220.5,
						},
						Current: types.ThreePhase{
							PhaseA: 3.5,
							PhaseB: 4.5,
							PhaseC: 5.5,
						},
						Frequency:         55,
						PowerFactor:       98,
						ActivePowerKW:     10,
						ReactivePowerKVAR: 20,
						ApparentPowerKVA:  30,
						EnergyConsumedKWH: 1234.5,
						EnergyExportedKWH: 5,
					},
					Events: []types.Event{},
					Location: types.Location{
						Latitude:  23.45663,
						Longitude: 12.436432,
					},
					SignalStrengthDBM:   -67,
					BatteryLevelPercent: 90,
					Status:              "OK",
				}

				payload, err := json.Marshal(data)
				if err != nil {
					log.Printf("error marshalling data: %v", err)
					continue
				}

				if token := client.Connect(); token.Wait() && token.Error() != nil {
					log.Fatalf("Failed to connect: %v", token.Error())
				}
				defer client.Disconnect(250)

				token := client.Publish("loggers", 0, false, payload)
				token.Wait()
			}
		}()
	}

	select {} // block forever
}
