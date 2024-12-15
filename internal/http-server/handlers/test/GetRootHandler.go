package test

import (
	"net/http"
)

// @Summary Тестовый эндпоинт
// @Description Это описание эндпоинта
// @Tags пример23
// @Success 200 {string} string "Техническая информация"
// @Router / [get]
func GetRootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Сервер запущен"))
	if err != nil {
		return
	}
}
