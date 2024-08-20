package main

import (
	"context"
	"fmt"
	"github.com/NikKazzzzzz/Scheduler/internal/config"
	event24 "github.com/NikKazzzzzz/Scheduler/internal/events"
	"github.com/NikKazzzzzz/Scheduler/internal/rabbitmq"
	"github.com/NikKazzzzzz/Scheduler/lib/sl"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting scheduler",
		slog.String("env", cfg.Env))

	var mongoDSN string
	if cfg.Env == envLocal {
		// Для локальной среды
		mongoDSN = fmt.Sprintf("mongodb://%s:%s@%s:27017/?authSource=admin", cfg.Database.Username, cfg.Database.Password, "localhost")
	} else {
		// Для Docker среды
		// Замена username и password в строке подключения
		username := cfg.Database.Username
		if username == "" {
			username = os.Getenv("MONGO_USERNAME")
		}

		password := cfg.Database.Password
		if password == "" {
			password = os.Getenv("MONGO_PASSWORD")
		}

		mongoDSN = strings.Replace(cfg.Database.MongoDSN, "username", username, 1)
		mongoDSN = strings.Replace(mongoDSN, "password", password, 1)
	}

	clientOptions := options.Client().ApplyURI(mongoDSN)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error("failed to connect to database: %v", sl.Err(err))
		os.Exit(1)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Error("failed to ping database: %v", sl.Err(err))
		os.Exit(1)
	}

	if cfg.Database.DatabaseName == "" {
		log.Error("database name cannot be empty")
		os.Exit(1)
	}

	collection := client.Database(cfg.Database.DatabaseName).Collection("events")
	eventService := event24.NewEventService(collection)

	producer, err := rabbitmq.NewProducer(cfg.RabbitMQ.URL, cfg.RabbitMQ.Queue)
	if err != nil {
		log.Error("failed to create producer: %v", sl.Err(err))
	}
	defer producer.Channel.Close()

	log.Info("Scheduler server started and is running...")

	for {
		events, err := eventService.GetEventsInNext24Hours()
		if err != nil {
			log.Error("failed to get events: %v", sl.Err(err))
			continue
		}

		for _, event := range events {
			body := event.Title
			log.Info("Preparing to publish event: " + body)

			err := producer.PublishEvent(body)
			if err != nil {
				log.Error("failed to publish event: %v", sl.Err(err))
			} else {
				log.Info("Successfully published event: " + body)
			}
		}

		time.Sleep(cfg.Scheduler.CheckInterval)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return log
}
