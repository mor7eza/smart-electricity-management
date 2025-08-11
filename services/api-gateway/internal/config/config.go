package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server Server `mapstructure:"server"`
}

type Server struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

func LoadConfig() (Config, error) {
	viper.AddConfigPath("./internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
