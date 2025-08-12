package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MQTT MQTT `mapstructure:"mqtt"`
}

type MQTT struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

func LoadConfig() Config {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config data: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("error unmarshalling config data: %v", err)
	}

	return config
}
