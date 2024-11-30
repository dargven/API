package bookingrepository

import (
	"API/internal/Storage/postrgeSQL"
	"API/internal/lib/logger/sl"
	book "API/internal/models/booking"
	"context"
	"fmt"
	"log/slog"
	"os"
)

type bookingRep struct {
	h      *postrgeSQL.Database
	logger *slog.Logger
}

func NewBookingRep(db *postrgeSQL.Database) *bookingRep {
	return &bookingRep{
		h:      db,
		logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
}

func (h *bookingRep) bookingMigrations() error {
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
		h.logger.Error("Failed to create migration for table booking", sl.Err(err))
		return fmt.Errorf("failed to create migration for table booking: %w", err)
	}
	h.logger.Info("Migration completed successfully (table created or already exists)")
	return nil

}

func (h *bookingRep) AddBooking(book book.Booking) (int64, error) {
	query := `INSERT INTO bookings (user_id, event_id, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := h.h.Pool.QueryRow(context.Background(), query, book.User_id, book.Event_id, book.Status, book.Created_at).Scan(&book.ID)
	if err != nil {
		h.logger.Error("failed to create booking")

		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return book.ID, nil
}

func (h *bookingRep) GetBookings(eventID int) ([]book.Booking, error) { // пока так
	query := `SELECT id, user_id, event_id, status, created_at FROM bookings WHERE event_id = $1;`

	rows, err := h.h.Pool.Query(context.Background(), query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}
	defer rows.Close()

	var bookings []book.Booking
	for rows.Next() {
		var booking book.Booking
		if err := rows.Scan(&booking.ID, &booking.User_id, &booking.Event_id, &booking.Status, &booking.Created_at); err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return bookings, nil
}

func (h *bookingRep) PostBooking(book book.Booking) error {
	query := `UPDATE bookings SET status = $1 WHERE id = $2;`
	_, err := h.h.Pool.Exec(context.Background(), query, book.Status, book.ID)
	if err != nil {
		h.logger.Error("failed to change status")

		return fmt.Errorf("failed to change status: %w", err)
	}

	return nil
}

func (h *bookingRep) DelBooking(book book.Booking) error {
	query := `DELETE FROM bookings WHERE id = $1;`
	_, err := h.h.Pool.Exec(context.Background(), query, book.ID)
	if err != nil {
		h.logger.Error("failed to delete booking")

		return fmt.Errorf("failed to delete booking: %w", err)
	}
	return nil
}
