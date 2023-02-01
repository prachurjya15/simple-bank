package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DbDriver      string        `mapstructure:"DB_DRIVER"`
	DbSource      string        `mapstructure:"DB_SOURCE"`
	Address       string        `mapstructure:"ADDRESS"`
	SymmetricKey  string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
