package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	UpdateTime time.Duration
	APIKey     string
}

func MustLoad() *Config {
	key, t := "API_KEY", "UPDATE_TIME"

	var c Config
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic("error reading a .env file")
	}
	if viper.IsSet(key) {
		c.APIKey = viper.GetString(key)
	} else {
		c.APIKey = ""
	}

	if viper.IsSet(t) {
		c.UpdateTime = viper.GetDuration(t)
	} else {
		c.UpdateTime = time.Second * 600
	}

	return &c
}
