package repository

import (
	"database/sql"
	"time"
)

type TrackingEvent struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"order_id"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	CreatedAt   time.Time `json:"created_at"`
}

func SaveEvent(db *sql.DB, e TrackingEvent) error {
	_, err := db.Exec(
		`INSERT INTO tracking_events (id, order_id, location, description, event_date, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		e.ID, e.OrderID, e.Location, e.Description, e.EventDate, e.CreatedAt,
	)
	return err
}

func GetEventsByOrderID(db *sql.DB, orderID string) ([]TrackingEvent, error) {
	rows, err := db.Query(
		`SELECT id, order_id, location, description, event_date, created_at
		 FROM tracking_events WHERE order_id = $1 ORDER BY event_date`, orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []TrackingEvent
	for rows.Next() {
		var e TrackingEvent
		if err := rows.Scan(&e.ID, &e.OrderID, &e.Location, &e.Description, &e.EventDate, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
