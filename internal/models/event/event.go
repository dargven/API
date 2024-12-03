package event

import (
	user "API/internal/models/user"
	"time"
)

type Status struct {
	Pending   string
	Confirmed string
	Cancelled string
}

type Event struct {
	ID          int64     `json:"id"`
	UserId      user.User `json:"user_id" validate:"required"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" validate:"required"`
	Location    string    `json:"location" validate:"required"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type Response struct {
	Token     string
	EventId   int64     `json:"event_id" validate:"required"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
