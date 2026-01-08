package models

import "time"

// Event представляет мероприятие
type Event struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Category         string    `json:"category"` // категория (концерт, спорт, театр)
	ImageURL         *string   `json:"image_url,omitempty"`
	Venue            string    `json:"venue"`             // место проведения
	Address          string    `json:"address"`           // адрес
	Price            float64   `json:"price"`             // цена билета
	Capacity         int       `json:"capacity"`          // вместимость
	AvailableTickets int       `json:"available_tickets"` // доступные билеты
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	CreatorID        int64     `json:"creator_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// EventResponse - DTO для ответа
type EventResponse struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Category         string    `json:"category"`
	ImageURL         *string   `json:"image_url,omitempty"`
	Venue            string    `json:"venue"`
	Address          string    `json:"address"`
	Price            float64   `json:"price"`
	Capacity         int       `json:"capacity"`
	AvailableTickets int       `json:"available_tickets"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	CreatorID        int64     `json:"creator_id"`
	CreatedAt        time.Time `json:"created_at"`
}

// ToResponse конвертирует Event в EventResponse
func (e *Event) ToResponse() EventResponse {
	return EventResponse{
		ID:               e.ID,
		Title:            e.Title,
		Description:      e.Description,
		Category:         e.Category,
		ImageURL:         e.ImageURL,
		Venue:            e.Venue,
		Address:          e.Address,
		Price:            e.Price,
		Capacity:         e.Capacity,
		AvailableTickets: e.AvailableTickets,
		StartTime:        e.StartTime,
		EndTime:          e.EndTime,
		CreatorID:        e.CreatorID,
		CreatedAt:        e.CreatedAt,
	}
}
