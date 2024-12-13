package eventRepository

import (
	"API/internal/Storage/postrgeSQL"
	"API/internal/lib/logger/sl"
	event "API/internal/models/event"
	"context"
	"fmt"
	"log/slog"
	"os"
)

type EventRep struct {
	h      *postrgeSQL.Database
	logger *slog.Logger
}

func NewEventRep(db *postrgeSQL.Database) *EventRep {
	return &EventRep{
		h:      db,
		logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
}

func (h *EventRep) EventMigrations() error {
	query :=
		` 
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		title TEXT, 
		description TEXT, 
		date TIMESTAMP DEFAULT NOW(), 
		location TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);
	`
	_, err := h.h.Pool.Exec(context.Background(), query)
	if err != nil {
		h.logger.Error("Failed to create migration for table event", sl.Err(err))
		return fmt.Errorf("failed to create migration for table event: %w", err)
	}
	h.logger.Info("Migration completed successfully (table created or already exists)")
	return nil

}

func (h *EventRep) IsTitleUnique(event event.Event) (bool, error) {
	query := `SELECT NOT EXISTS (SELECT 1 FROM events WHERE title = $1)`
	var isUnique bool
	err := h.h.Pool.QueryRow(context.Background(), query, event.Title).Scan(&isUnique)
	if err != nil {
		h.logger.Error("failed to check title uniqueness ", sl.Err(err))
		return false, fmt.Errorf(" failed to check title uniqueness: %w", err)
	}

	return isUnique, nil
}

func (h *EventRep) IsEventExist(event event.Event) (bool, error) {
	query := `SELECT NOT EXISTS (SELECT 1 FROM events WHERE id = $1)`
	var isExist bool
	err := h.h.Pool.QueryRow(context.Background(), query, event.Title).Scan(&isExist)
	if err != nil {
		h.logger.Error("failed to check existion of this event ", sl.Err(err))
		return false, fmt.Errorf(" failed to check existion of this event:  %w", err)
	}

	return isExist, nil
}

func (h *EventRep) AddEvent(event event.Event) (int64, error) {
	query := `INSERT INTO events (user_id, title, description, date, location, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	unique, err := h.IsTitleUnique(event)

	if err != nil {
		return 0, fmt.Errorf("failed to check title uniqueness: %w", err)
	}

	if unique {

		err := h.h.Pool.QueryRow(context.Background(), query, event.UserId, event.Title, event.Description, event.Date, event.Location, event.Status, event.CreatedAt).Scan(&event.ID)
		if err != nil {
			h.logger.Error("failed to create event")

			return 0, fmt.Errorf("failed to create event: %w", err)
		}

	} else if !unique {
		h.logger.Error("event title is not unique")

		return 0, fmt.Errorf("event title is not unique")
	}

	return event.ID, nil
}

func (h *EventRep) GetEvent(id int) (event.Event, error) {
	query := `SELECT id, user_id, title, description, date, location, status, created_at FROM events WHERE id = $1;`

	var ev event.Event

	exist, err := h.IsEventExist(ev)

	if err != nil {
		return ev, fmt.Errorf("failed to check title uniqueness: %w", err)
	}

	if exist {
		err := h.h.Pool.QueryRow(context.Background(), query, id).Scan(&ev.ID, &ev.UserId, &ev.Title, &ev.Description, &ev.Date, &ev.Location, &ev.Status, &ev.CreatedAt)
		if err != nil {
			h.logger.Error("event with ID not found: ", id, " || ", sl.Err(err))
			return ev, fmt.Errorf("event with ID %d not found", id)
		}
	} else if !exist {
		h.logger.Error("event with that id is not exist", sl.Err(err))
		return ev, fmt.Errorf("event with that id is not exist", err)
	}

	return ev, nil
}

func (h *EventRep) GetAllEvents() ([]event.Event, error) {
	query := `SELECT * FROM events`

	rows, err := h.h.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of events: %w", err)
	}
	defer rows.Close()

	events := make([]event.Event, 0)

	for rows.Next() {
		var ev event.Event
		if err := rows.Scan(&ev.ID, &ev.UserId, &ev.Title, &ev.Description, &ev.Date, &ev.Location, &ev.Status, &ev.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, ev)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return events, nil
}

func (h *EventRep) PostEvent(event event.Event) error {
	query := `UPDATE events SET status = $1 WHERE id = $2;`
	_, err := h.h.Pool.Exec(context.Background(), query, event.Status, event.ID)
	if err != nil {
		h.logger.Error("failed to change status")

		return fmt.Errorf("failed to change status: %w", err)
	}

	return nil
}

//потенциал на изменение но пока так

func (h *EventRep) DelEvent(event event.Event) error {
	query := `DELETE FROM events WHERE id = $1;`
	_, err := h.h.Pool.Exec(context.Background(), query, event.ID)
	if err != nil {
		h.logger.Error("failed to delete event")

		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}
