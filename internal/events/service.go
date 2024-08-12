package events

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          int
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

type EventService struct {
	DB *sql.DB
}

func NewEventService(db *sql.DB) *EventService {
	return &EventService{DB: db}
}

func (s *EventService) GetEventsInNext24Hours() ([]Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time 
		FROM events 
		WHERE start_time BETWEEN NOW() AND NOW() + INTERVAL '1 day'
	`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var startTimeStr, endTimeStr string

		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &startTimeStr, &endTimeStr); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
