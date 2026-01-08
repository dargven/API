package models

import "time"

// Event представляет мероприятие
type Event struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatorID   int64     `json:"creator_id"`
	MaxSlots    int       `json:"max_slots"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EventResponse - DTO для ответа без чувствительных данных
type EventResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatorID   int64     `json:"creator_id"`
	MaxSlots    int       `json:"max_slots"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToResponse конвертирует Event в EventResponse
func (e *Event) ToResponse() EventResponse {
	return EventResponse{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Location:    e.Location,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		CreatorID:   e.CreatorID,
		MaxSlots:    e.MaxSlots,
		CreatedAt:   e.CreatedAt,
	}
}
