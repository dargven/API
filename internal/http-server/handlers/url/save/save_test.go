package save

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"

// 	"API/internal/lib/logger/sl"
// )

// func TestSaveHandler(t *testing.T) {
// 	cases := []struct {
// 		name      string // Имя теста
// 		alias     string // Отправляемый alias
// 		url       string // Отправляемый URL
// 		respError string //Указываем какую ошибку хотим получить
// 		mockError error  //Ошибка которую выдает mock
// 	}{
// 		{
// 			name:  "Success",
// 			alias: "test_alias",
// 			url:   "http://google.com",
// 			// Тут поля respError и mockError оставляем пустыми,
// 			// т.к. это успешный запрос
// 		}, {
// 			name:      "Empty url",
// 			alias:     "test_alias",
// 			url:       "",
// 			respError: "URL is required",
// 			mockError: nil,
// 		}, {
// 			name:      "Inbalid URL",
// 			alias:     "test_alias",
// 			url:       "invalid-url",
// 			respError: "invalid url format",
// 			mockError: nil,
// 		}, {
// 			name:      "Missing Alias",
// 			alias:     "",
// 			url:       "http://google.com",
// 			respError: "Alias is required",
// 			mockError: nil,
// 		}, {
// 			name:      "Duplicate Alias",
// 			alias:     "existing_alias",
// 			url:       "http://google.com",
// 			respError: "Alias already exists",
// 			mockError: errors.New("alias already exists"),
// 		}, {
// 			name:      "Save URL Error",
// 			alias:     "test_alias",
// 			url:       "http://google.com",
// 			respError: "Failed to save URL",
// 			mockError: errors.New("database connection error"),
// 		},
// 	}

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			//создаем объект мока стораджа
// 			urlSaverMock := mocks.NewURLSaver(t)
// 			//если ожидается успешный ответ значит к моку точно будет вызов
// 			//даже если ожидаем ошибку
// 			//но мок должен ответить ошибкой, к нему тоже будет запрос
// 			if tc.respError == "" || tc.mockError != nil {
// 				// Сообщаем моку, какой к нему будет запрос, и что надо вернуть
// 				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
// 					Return(int64(1), tc.mockError).
// 					Once() // Запрос будет ровно один
// 			}

// 			// Создаем наш хэндлер
// 			handler := New(sl.NewDiscardLogger(), urlSaverMock)

// 			// Формируем тело запроса
// 			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

// 			// Создаем объект запроса
// 			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
// 			require.NoError(t, err)

// 			// Создаем ResponseRecorder для записи ответа хэндлера
// 			rr := httptest.NewRecorder()
// 			// Обрабатываем запрос, записывая ответ в рекордер
// 			handler.ServeHTTP(rr, req)

// 			// Проверяем, что статус ответа корректный
// 			require.Equal(t, rr.Code, http.StatusOK)

// 			body := rr.Body.String()

// 			var resp Response

// 			// Анмаршаллим тело, и проверяем что при этом не возникло ошибок
// 			require.NoError(t, json.Unmarshal([]byte(body), &resp))

// 			// Проверяем наличие требуемой ошибки в ответе
// 			require.Equal(t, tc.respError, resp.Error)

// 			// Проверяем, что вызов SaveURL был один раз, если ожидался
// 			urlSaverMock.AssertExpectations(t)

// 		})
// 	}

// }
