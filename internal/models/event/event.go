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
	ID        int64     `json:"id"`
	UserId    user.User `json:"user_id" validate:"required"`
	EventId   int64     `json:"event_id" validate:"required"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	token     string
	EventId   int64     `json:"event_id" validate:"required"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
