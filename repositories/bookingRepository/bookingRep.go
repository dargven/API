package bookingrepository

import (
	"API/internal/Storage/postrgeSQL"
	"API/internal/lib/logger/sl"
	event "API/internal/models/event"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

type BookingRep struct {
	h      *postrgeSQL.Database
	logger *slog.Logger
}

func NewBookingRep(db *postrgeSQL.Database) *BookingRep {
	return &BookingRep{
		h:      db,
		logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
}

func (h *BookingRep) bookingMigrations() error {
	query :=
		` 
	CREATE TABLE IF NOT EXISTS bookings (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		event_id INT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE
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

func (h *BookingRep) AddBooking(event event.Event) (int64, error) {
	query := `INSERT INTO bookings (user_id, event_id, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := h.h.Pool.QueryRow(context.Background(), query, event.UserId, event.EventId, event.Status, event.CreatedAt).Scan(&event.ID)
	if err != nil {
		h.logger.Error("failed to create event")

		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return event.ID, nil
}

func (h *BookingRep) GetBookings(eventID int) ([]event.Event, error) { // пока так
	query := `SELECT id, user_id, event_id, status, created_at FROM bookings WHERE event_id = $1;`

	rows, err := h.h.Pool.Query(context.Background(), query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}
	defer rows.Close()

	var bookings []event.Event
	for rows.Next() {
		var booking event.Event
		if err := rows.Scan(&booking.ID, &booking.UserId, &booking.EventId, &booking.Status, &booking.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return bookings, nil
}

func (h *BookingRep) PostBooking(book event.Event) error {
	query := `UPDATE bookings SET status = $1 WHERE id = $2;`
	_, err := h.h.Pool.Exec(context.Background(), query, book.Status, book.ID)
	if err != nil {
		h.logger.Error("failed to change status")

		return fmt.Errorf("failed to change status: %w", err)
	}

	return nil
}

func (h *BookingRep) DelBooking(book event.Event) error {
	query := `DELETE FROM bookings WHERE id = $1;`
	_, err := h.h.Pool.Exec(context.Background(), query, book.ID)
	if err != nil {
		h.logger.Error("failed to delete event")

		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}
func (h *BookingRep) FetchEventByID(eventID int) (event.Event, error) {
	query := "SELECT id, title, description, date, location FROM events WHERE id = $1"
	var fetchedEvent event.Event

	// Выполнение запроса
	err := h.h.Pool.QueryRow(context.Background(), query, eventID).Scan(
		&fetchedEvent.ID,
		//&fetchedEvent.Title,
		//&fetchedEvent.Description,
		//&fetchedEvent.Date,
		//&fetchedEvent.Location,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fetchedEvent, fmt.Errorf("event with id %d not found", eventID)
		}
		return fetchedEvent, fmt.Errorf("failed to fetch event: %w", err)
	}

	return fetchedEvent, nil
}
