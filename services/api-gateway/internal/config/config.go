package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Server Server `mapstructure:"server"`
}

type Server struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

func LoadConfig() Config {
	viper.AddConfigPath("./internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("error reading config: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logrus.Fatalf("error unmarshalling config data: %v", err)
	}

	return config
}
