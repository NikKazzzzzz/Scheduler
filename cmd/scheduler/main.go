package main

import (
	"database/sql"
	"github.com/NikKazzzzzz/Scheduler/internal/config"
	event24 "github.com/NikKazzzzzz/Scheduler/internal/events"
	"github.com/NikKazzzzzz/Scheduler/internal/rabbitmq"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("./config/scheduler.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	eventService := event24.NewEventService(db)

	producer, err := rabbitmq.NewProducer(cfg.RabbitMQ.URL, cfg.RabbitMQ.Queue)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	defer producer.Channel.Close()

	log.Println("Scheduler server started and is running...")

	for {
		events, err := eventService.GetEventsInNext24Hours()
		if err != nil {
			log.Printf("failed to get events: %v", err)
			continue
		}

		for _, event := range events {
			body := event.Title
			log.Printf("Preparing to publish event: %s", body)

			err := producer.PublishEvent(body)
			if err != nil {
				log.Printf("failed to publish event: %v", err)
			} else {
				log.Printf("Successfully published event: %s", body)
			}
		}

		time.Sleep(cfg.Scheduler.CheckInterval)
	}
}
