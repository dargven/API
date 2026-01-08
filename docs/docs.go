// Package docs API для бронирования мероприятий
//
// API для регистрации, авторизации и управления мероприятиями
//
//	Schemes: http, https
//	Host: localhost:8082
//	BasePath: /
//	Version: 2.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	SecurityDefinitions:
//	  Bearer:
//	    type: apiKey
//	    name: Authorization
//	    in: header
//	    description: JWT токен в формате "Bearer {token}"
//
// swagger:meta
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "API для бронирования мероприятий с JWT авторизацией. Включает регистрацию, профиль пользователя, создание мероприятий, бронирование билетов и полнотекстовый поиск.",
        "title": "Event Booking API",
        "contact": {
            "name": "API Support",
            "email": "support@example.com"
        },
        "version": "2.0.0"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/register": {
            "post": {
                "description": "Создает нового пользователя и возвращает JWT токен",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["auth"],
                "summary": "Регистрация пользователя",
                "parameters": [
                    {
                        "description": "Данные для регистрации",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Пользователь создан",
                        "schema": {
                            "$ref": "#/definitions/RegisterResponse"
                        }
                    },
                    "400": {
                        "description": "Неверные данные",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "409": {
                        "description": "Email уже существует",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Авторизует пользователя и возвращает JWT токен",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["auth"],
                "summary": "Авторизация пользователя",
                "parameters": [
                    {
                        "description": "Данные для авторизации",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная авторизация",
                        "schema": {
                            "$ref": "#/definitions/LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Неверные данные",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "401": {
                        "description": "Неверный email или пароль",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/profile": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Возвращает профиль текущего пользователя с балансом",
                "produces": ["application/json"],
                "tags": ["profile"],
                "summary": "Получить профиль",
                "responses": {
                    "200": {
                        "description": "Профиль пользователя",
                        "schema": {
                            "$ref": "#/definitions/ProfileResponse"
                        }
                    },
                    "401": {
                        "description": "Не авторизован",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            },
            "put": {
                "security": [{"Bearer": []}],
                "description": "Обновляет профиль текущего пользователя",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["profile"],
                "summary": "Обновить профиль",
                "parameters": [
                    {
                        "description": "Данные профиля",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdateProfileRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Профиль обновлен",
                        "schema": {
                            "$ref": "#/definitions/ProfileResponse"
                        }
                    },
                    "400": {
                        "description": "Неверные данные",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "401": {
                        "description": "Не авторизован",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/profile/balance": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "Пополняет баланс текущего пользователя",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["profile"],
                "summary": "Пополнить баланс",
                "parameters": [
                    {
                        "description": "Сумма пополнения",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/TopUpBalanceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Баланс пополнен",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "400": {
                        "description": "Неверные данные",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "401": {
                        "description": "Не авторизован",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/events": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Возвращает список всех мероприятий с пагинацией",
                "produces": ["application/json"],
                "tags": ["events"],
                "summary": "Получение списка мероприятий",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Количество записей (по умолчанию 20, максимум 100)",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Смещение (по умолчанию 0)",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список мероприятий",
                        "schema": {
                            "$ref": "#/definitions/GetAllEventsResponse"
                        }
                    },
                    "401": {
                        "description": "Не авторизован",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "Создает новое мероприятие с ценой и количеством билетов",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["events"],
                "summary": "Создание мероприятия",
                "parameters": [
                    {
                        "description": "Данные мероприятия",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreateEventRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Мероприятие создано",
                        "schema": {
                            "$ref": "#/definitions/CreateEventResponse"
                        }
                    },
                    "400": {
                        "description": "Неверные данные",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "401": {
                        "description": "Не авторизован",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/events/{id}": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Возвращает мероприятие по его ID",
                "produces": ["application/json"],
                "tags": ["events"],
                "summary": "Получение мероприятия по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID мероприятия",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Мероприятие",
                        "schema": {
                            "$ref": "#/definitions/GetEventResponse"
                        }
                    },
                    "401": {
                        "description": "Не авторизован",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "404": {
                        "description": "Мероприятие не найдено",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/events/{id}/book": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "Бронирует билет на мероприятие. Списывает деньги с баланса пользователя.",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["bookings"],
                "summary": "Забронировать билет",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID мероприятия",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Количество билетов",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreateBookingRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Бронирование создано",
                        "schema": {
                            "$ref": "#/definitions/CreateBookingResponse"
                        }
                    },
                    "400": {
                        "description": "Неверные данные",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "401": {
                        "description": "Не авторизован",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "404": {
                        "description": "Мероприятие не найдено",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "409": {
                        "description": "Бронирование уже существует",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "422": {
                        "description": "Недостаточно билетов или баланса",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/bookings": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Возвращает список всех бронирований текущего пользователя",
                "produces": ["application/json"],
                "tags": ["bookings"],
                "summary": "Мои билеты",
                "responses": {
                    "200": {
                        "description": "Список бронирований",
                        "schema": {
                            "$ref": "#/definitions/ListBookingsResponse"
                        }
                    },
                    "401": {
                        "description": "Не авторизован",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/bookings/{id}": {
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Отменяет бронирование и возвращает деньги на баланс",
                "produces": ["application/json"],
                "tags": ["bookings"],
                "summary": "Отменить бронь",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID бронирования",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Бронирование отменено",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "401": {
                        "description": "Не авторизован",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "404": {
                        "description": "Бронирование не найдено",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/search": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Полнотекстовый поиск мероприятий с фильтрами по категории, дате и цене",
                "produces": ["application/json"],
                "tags": ["search"],
                "summary": "Поиск мероприятий",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Поисковый запрос",
                        "name": "q",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Категория (concert, sport, theater, exhibition, festival, other)",
                        "name": "category",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Дата от (RFC3339)",
                        "name": "date_from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Дата до (RFC3339)",
                        "name": "date_to",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "Минимальная цена",
                        "name": "price_min",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "Максимальная цена",
                        "name": "price_max",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Лимит (по умолчанию 20, максимум 100)",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Смещение",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Результаты поиска",
                        "schema": {
                            "$ref": "#/definitions/SearchResponse"
                        }
                    },
                    "400": {
                        "description": "Неверные параметры",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "401": {
                        "description": "Не авторизован",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Response": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "error": {
                    "type": "string",
                    "example": "error message"
                }
            }
        },
        "RegisterRequest": {
            "type": "object",
            "required": ["email", "name", "password"],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "password": {
                    "type": "string",
                    "example": "securePassword123"
                }
            }
        },
        "RegisterResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "user": {
                    "$ref": "#/definitions/UserResponse"
                },
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIs..."
                }
            }
        },
        "LoginRequest": {
            "type": "object",
            "required": ["email", "password"],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "password": {
                    "type": "string",
                    "example": "securePassword123"
                }
            }
        },
        "LoginResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIs..."
                }
            }
        },
        "UserResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "created_at": {
                    "type": "string",
                    "example": "2024-01-15T10:00:00Z"
                }
            }
        },
        "ProfileResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "profile": {
                    "$ref": "#/definitions/Profile"
                }
            }
        },
        "Profile": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "phone": {
                    "type": "string",
                    "example": "+79001234567"
                },
                "avatar_url": {
                    "type": "string",
                    "example": "https://example.com/avatar.jpg"
                },
                "bio": {
                    "type": "string",
                    "example": "Software Developer"
                },
                "balance": {
                    "type": "number",
                    "example": 5000.00
                },
                "created_at": {
                    "type": "string",
                    "example": "2024-01-15T10:00:00Z"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2024-01-15T10:00:00Z"
                }
            }
        },
        "UpdateProfileRequest": {
            "type": "object",
            "required": ["name"],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "John Updated"
                },
                "phone": {
                    "type": "string",
                    "example": "+79001234567"
                },
                "avatar_url": {
                    "type": "string",
                    "example": "https://example.com/avatar.jpg"
                },
                "bio": {
                    "type": "string",
                    "example": "Senior Developer"
                }
            }
        },
        "TopUpBalanceRequest": {
            "type": "object",
            "required": ["amount"],
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 1000.00
                }
            }
        },
        "CreateEventRequest": {
            "type": "object",
            "required": ["title", "category", "venue", "address", "capacity", "start_time", "end_time"],
            "properties": {
                "title": {
                    "type": "string",
                    "example": "Рок-концерт 2024"
                },
                "description": {
                    "type": "string",
                    "example": "Лучший рок-концерт года!"
                },
                "category": {
                    "type": "string",
                    "enum": ["concert", "sport", "theater", "exhibition", "festival", "other"],
                    "example": "concert"
                },
                "image_url": {
                    "type": "string",
                    "example": "https://example.com/concert.jpg"
                },
                "venue": {
                    "type": "string",
                    "example": "Олимпийский"
                },
                "address": {
                    "type": "string",
                    "example": "Москва, Олимпийский проспект, 16"
                },
                "price": {
                    "type": "number",
                    "example": 3500.00
                },
                "capacity": {
                    "type": "integer",
                    "example": 15000
                },
                "start_time": {
                    "type": "string",
                    "example": "2024-12-20T19:00:00Z"
                },
                "end_time": {
                    "type": "string",
                    "example": "2024-12-20T23:00:00Z"
                }
            }
        },
        "CreateEventResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "event": {
                    "$ref": "#/definitions/EventResponse"
                }
            }
        },
        "GetEventResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "event": {
                    "$ref": "#/definitions/EventResponse"
                }
            }
        },
        "GetAllEventsResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/EventResponse"
                    }
                },
                "total": {
                    "type": "integer",
                    "example": 10
                }
            }
        },
        "EventResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "title": {
                    "type": "string",
                    "example": "Рок-концерт 2024"
                },
                "description": {
                    "type": "string",
                    "example": "Лучший рок-концерт года!"
                },
                "category": {
                    "type": "string",
                    "example": "concert"
                },
                "image_url": {
                    "type": "string",
                    "example": "https://example.com/concert.jpg"
                },
                "venue": {
                    "type": "string",
                    "example": "Олимпийский"
                },
                "address": {
                    "type": "string",
                    "example": "Москва, Олимпийский проспект, 16"
                },
                "price": {
                    "type": "number",
                    "example": 3500.00
                },
                "capacity": {
                    "type": "integer",
                    "example": 15000
                },
                "available_tickets": {
                    "type": "integer",
                    "example": 14500
                },
                "start_time": {
                    "type": "string",
                    "example": "2024-12-20T19:00:00Z"
                },
                "end_time": {
                    "type": "string",
                    "example": "2024-12-20T23:00:00Z"
                },
                "creator_id": {
                    "type": "integer",
                    "example": 1
                },
                "created_at": {
                    "type": "string",
                    "example": "2024-01-15T10:00:00Z"
                }
            }
        },
        "CreateBookingRequest": {
            "type": "object",
            "required": ["quantity"],
            "properties": {
                "quantity": {
                    "type": "integer",
                    "minimum": 1,
                    "maximum": 10,
                    "example": 2
                }
            }
        },
        "CreateBookingResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "booking": {
                    "$ref": "#/definitions/Booking"
                }
            }
        },
        "ListBookingsResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "bookings": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/BookingResponse"
                    }
                }
            }
        },
        "Booking": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                },
                "event_id": {
                    "type": "integer",
                    "example": 1
                },
                "quantity": {
                    "type": "integer",
                    "example": 2
                },
                "total_price": {
                    "type": "number",
                    "example": 7000.00
                },
                "status": {
                    "type": "string",
                    "enum": ["confirmed", "cancelled", "used"],
                    "example": "confirmed"
                },
                "booking_code": {
                    "type": "string",
                    "example": "BK-a1b2c3d4e5f6g7h8"
                },
                "created_at": {
                    "type": "string",
                    "example": "2024-01-15T10:00:00Z"
                }
            }
        },
        "BookingResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "event_id": {
                    "type": "integer",
                    "example": 1
                },
                "event_title": {
                    "type": "string",
                    "example": "Рок-концерт 2024"
                },
                "event_date": {
                    "type": "string",
                    "example": "2024-12-20T19:00:00Z"
                },
                "venue": {
                    "type": "string",
                    "example": "Олимпийский"
                },
                "quantity": {
                    "type": "integer",
                    "example": 2
                },
                "total_price": {
                    "type": "number",
                    "example": 7000.00
                },
                "status": {
                    "type": "string",
                    "example": "confirmed"
                },
                "booking_code": {
                    "type": "string",
                    "example": "BK-a1b2c3d4e5f6g7h8"
                },
                "created_at": {
                    "type": "string",
                    "example": "2024-01-15T10:00:00Z"
                }
            }
        },
        "SearchResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/EventResponse"
                    }
                },
                "total": {
                    "type": "integer",
                    "example": 50
                },
                "limit": {
                    "type": "integer",
                    "example": 20
                },
                "offset": {
                    "type": "integer",
                    "example": 0
                },
                "has_more": {
                    "type": "boolean",
                    "example": true
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "JWT токен. Формат: Bearer {token}"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "2.0.0",
	Host:             "localhost:8082",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Event Booking API",
	Description:      "API для бронирования мероприятий с JWT авторизацией",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
