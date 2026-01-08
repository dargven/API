package models

import "time"

// BookingStatus статусы бронирования
type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusUsed      BookingStatus = "used"
)

// Booking представляет бронирование билета
type Booking struct {
	ID          int64         `json:"id"`
	UserID      int64         `json:"user_id"`
	EventID     int64         `json:"event_id"`
	Quantity    int           `json:"quantity"`
	TotalPrice  float64       `json:"total_price"`
	Status      BookingStatus `json:"status"`
	BookingCode string        `json:"booking_code"`
	CreatedAt   time.Time     `json:"created_at"`
}

// BookingResponse - DTO для ответа с данными мероприятия
type BookingResponse struct {
	ID          int64         `json:"id"`
	EventID     int64         `json:"event_id"`
	EventTitle  string        `json:"event_title"`
	EventDate   time.Time     `json:"event_date"`
	Venue       string        `json:"venue"`
	Quantity    int           `json:"quantity"`
	TotalPrice  float64       `json:"total_price"`
	Status      BookingStatus `json:"status"`
	BookingCode string        `json:"booking_code"`
	CreatedAt   time.Time     `json:"created_at"`
}

// BookingWithEvent - бронирование с информацией о мероприятии
type BookingWithEvent struct {
	Booking
	EventTitle string    `json:"event_title"`
	EventDate  time.Time `json:"event_date"`
	Venue      string    `json:"venue"`
}

// ToResponse конвертирует BookingWithEvent в BookingResponse
func (b *BookingWithEvent) ToResponse() BookingResponse {
	return BookingResponse{
		ID:          b.ID,
		EventID:     b.EventID,
		EventTitle:  b.EventTitle,
		EventDate:   b.EventDate,
		Venue:       b.Venue,
		Quantity:    b.Quantity,
		TotalPrice:  b.TotalPrice,
		Status:      b.Status,
		BookingCode: b.BookingCode,
		CreatedAt:   b.CreatedAt,
	}
}
