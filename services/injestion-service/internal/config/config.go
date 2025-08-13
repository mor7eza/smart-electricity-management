package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	MQTT MQTT `mapstructure:"mqtt"`
}

type MQTT struct {
	Address  string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Topic    string `mapstructure:"topic"`
	ClientID string `mapstructure:"client_id"`
}

func LoadConfig() Config {
	viper.AddConfigPath("./internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logrus.Fatalf("error unmarshalling config data: %v", err)
	}

	fmt.Println(config)
	return config
}
