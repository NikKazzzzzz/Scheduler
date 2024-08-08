package main

import (
	"database/sql"
	"github.com/NikKazzzzzz/Scheduler/internal/config"
	"github.com/NikKazzzzzz/Scheduler/internal/events"
	"github.com/NikKazzzzzz/Scheduler/internal/rabbirmq"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"log"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("./config/scheduler.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	eventService := events.NewEventService(db)

	producer, err := rabbirmq.NewProducer(cfg.RabbitMQ.URL, cfg.RabbitMQ.Queue)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	defer producer.Channel.Close()

	for {
		events, err := eventService.GetEventsInNext24Hours()
		if err != nil {
			log.Printf("failed to get events: %v", err)
			continue
		}

		for _, event := range events {
			body := event.Title
			err := producer.PublishEvent(body)
			if err != nil {
				log.Printf("failed to publish event: %v", err)
			}
		}

		time.Sleep(cfg.Scheduler.CheckInterval)
	}
}
