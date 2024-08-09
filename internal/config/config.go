package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	RabbitMQ struct {
		URL   string `mapstructure:"url"`
		Queue string `mapstructure:"queue"`
	} `mapstructure:"rabbitmq"`

	Database struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"database"`

	Scheduler struct {
		CheckInterval time.Duration `mapstructure:"check_interval"`
	} `mapstructure:"scheduler"`
}

func LoadConfig(dsn string) (*Config, error) {
	viper.SetConfigFile(dsn)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
