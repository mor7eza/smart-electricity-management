package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MQTT MQTT
}

type MQTT struct {
	Address  string
	Port     int
	Topic    string
	ClientID string
}

func LoadConfig() Config {
	viper.AutomaticEnv()

	config := Config{
		MQTT: MQTT{
			Address:  viper.GetString("MQTT_ADDRESS"),
			Port:     viper.GetInt("MQTT_PORT"),
			Topic:    viper.GetString("MQTT_TOPIC"),
			ClientID: viper.GetString("MQTT_CLIENT_ID"),
		},
	}

	return config
}
