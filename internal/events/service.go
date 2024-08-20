package events

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Event struct {
	ID          int       `bson:"id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	StartTime   time.Time `bson:"start_time"`
	EndTime     time.Time `bson:"end_time"`
}

type EventService struct {
	Collection *mongo.Collection
}

func NewEventService(collection *mongo.Collection) *EventService {
	return &EventService{Collection: collection}
}

func (s *EventService) GetEventsInNext24Hours() ([]Event, error) {
	now := time.Now()
	oneDayLater := now.Add(24 * time.Hour)

	filter := bson.M{
		"start_time": bson.M{"$gte": now},
		"end_time":   bson.M{"$lte": oneDayLater},
	}

	cursor, err := s.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var events []Event
	for cursor.Next(context.TODO()) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
