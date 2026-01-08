package slogdiscard

// Нужен Logger Для корректной работы тестов.
import (
	"context"

	"golang.org/x/exp/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	//Возвращает тот же обработчик, т.к нет атрибутов для сохранения
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// Возвращает тот же обработчик, т.к нет группы для сохранения

	return h
}
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	//Всегда возвращает false, т.к запись журнала игнорируется
	return false
}
