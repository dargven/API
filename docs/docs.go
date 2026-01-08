// Package docs API для бронирования мероприятий
//
// API для регистрации, авторизации и управления мероприятиями
//
//	Schemes: http, https
//	Host: localhost:8082
//	BasePath: /
//	Version: 1.0.0
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
        "description": "API для бронирования мероприятий с JWT авторизацией",
        "title": "Event Booking API",
        "contact": {
            "name": "API Support",
            "email": "support@example.com"
        },
        "version": "1.0.0"
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
                "description": "Создает новое мероприятие",
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
        "CreateEventRequest": {
            "type": "object",
            "required": ["title", "location", "start_time", "end_time", "max_slots"],
            "properties": {
                "title": {
                    "type": "string",
                    "example": "Конференция Go"
                },
                "description": {
                    "type": "string",
                    "example": "Ежегодная конференция разработчиков Go"
                },
                "location": {
                    "type": "string",
                    "example": "Москва, ул. Примерная 1"
                },
                "start_time": {
                    "type": "string",
                    "example": "2024-06-15T10:00:00Z"
                },
                "end_time": {
                    "type": "string",
                    "example": "2024-06-15T18:00:00Z"
                },
                "max_slots": {
                    "type": "integer",
                    "example": 100
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
                    "example": "Конференция Go"
                },
                "description": {
                    "type": "string",
                    "example": "Ежегодная конференция разработчиков Go"
                },
                "location": {
                    "type": "string",
                    "example": "Москва, ул. Примерная 1"
                },
                "start_time": {
                    "type": "string",
                    "example": "2024-06-15T10:00:00Z"
                },
                "end_time": {
                    "type": "string",
                    "example": "2024-06-15T18:00:00Z"
                },
                "creator_id": {
                    "type": "integer",
                    "example": 1
                },
                "max_slots": {
                    "type": "integer",
                    "example": 100
                },
                "booked_slots": {
                    "type": "integer",
                    "example": 25
                },
                "created_at": {
                    "type": "string",
                    "example": "2024-01-15T10:00:00Z"
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
	Version:          "1.0.0",
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
