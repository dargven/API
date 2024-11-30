package models

import (
	user "API/internal/models/user"
	"time"
)

type Status struct {
	Pending   string
	Confirmed string
	Cancelled string
}

type Booking struct {
	ID         int64     `json:"id"`
	User_id    user.User `json:"user_id" validate:"required"`
	Event_id   int64     `json:"event_id" validate:"required"`
	Status     Status    `json:"status"`
	Created_at time.Time `json:"created_at"`
}
