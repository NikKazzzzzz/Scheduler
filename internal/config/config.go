package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env string `yaml:"log_level" env-default:"local"`

	RabbitMQ struct {
		URL   string `yaml:"url"`
		Queue string `yaml:"queue"`
	} `yaml:"rabbitmq"`

	Database struct {
		MongoDSN     string `yaml:"mongo_dsn"`
		DatabaseName string `yaml:"databaseName"`
		Username     string `yaml:"username" env:"MONGO_USERNAME"`
		Password     string `yaml:"password" env:"MONGO_PASSWORD"`
	} `yaml:"database"`

	Scheduler struct {
		CheckInterval time.Duration `yaml:"check_interval"`
	} `yaml:"scheduler"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	log.Printf("loaded config: %+v", cfg)

	return &cfg
}
