# Миграции базы данных (PostgreSQL)

Этот документ содержит SQL скрипты для настройки базы данных PostgreSQL.

## Создание базы данных

```sql
-- Создание базы данных (выполнить от имени суперпользователя)
CREATE DATABASE events_db;

-- Подключение к базе данных
\c events_db
```

---

## Миграция 001: Таблица url

```sql
-- 001_create_url_table.sql
CREATE TABLE IF NOT EXISTS url (
    id BIGSERIAL PRIMARY KEY,
    alias VARCHAR(255) NOT NULL UNIQUE,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_url_alias ON url(alias);
```

---

## Миграция 002: Таблица users

```sql
-- 002_create_users_table.sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггер для автоматического обновления updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

---

## Миграция 003: Таблица events (мероприятия)

```sql
-- 003_create_events_table.sql
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    location VARCHAR(500) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    max_slots INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_events_creator_id ON events(creator_id);
CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);

-- Триггер для автоматического обновления updated_at
CREATE TRIGGER update_events_updated_at 
    BEFORE UPDATE ON events
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### Описание полей таблицы events

| Поле        | Тип                      | Описание                                |
|-------------|--------------------------|----------------------------------------|
| id          | BIGSERIAL                | Уникальный идентификатор               |
| title       | VARCHAR(200)             | Название мероприятия                   |
| description | TEXT                     | Описание мероприятия                   |
| location    | VARCHAR(500)             | Место проведения                       |
| start_time  | TIMESTAMP WITH TIME ZONE | Время начала                           |
| end_time    | TIMESTAMP WITH TIME ZONE | Время окончания                        |
| creator_id  | BIGINT                   | ID создателя (ссылка на users)         |
| max_slots   | INTEGER                  | Максимальное количество участников     |
| created_at  | TIMESTAMP WITH TIME ZONE | Дата создания записи                   |
| updated_at  | TIMESTAMP WITH TIME ZONE | Дата последнего обновления             |

---

## Будущие миграции

### Миграция 004: Таблица бронирований (опционально)

```sql
-- 004_create_bookings_table.sql
CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Уникальность: один пользователь - одна бронь на мероприятие
    UNIQUE(event_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_bookings_event_id ON bookings(event_id);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
```

---

## Применение миграций

### Через psql

```bash
psql -h localhost -U postgres -d events_db -f migrations/001_create_url_table.sql
psql -h localhost -U postgres -d events_db -f migrations/002_create_users_table.sql
psql -h localhost -U postgres -d events_db -f migrations/003_create_events_table.sql
```

---

## Откат миграций

### Откат миграции 003

```sql
-- rollback_003_create_events_table.sql
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP TABLE IF EXISTS events;
```

### Откат миграции 002

```sql
-- rollback_002_create_users_table.sql
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS users;
```

### Откат миграции 001

```sql
-- rollback_001_create_url_table.sql
DROP TABLE IF EXISTS url;
```

---

## Docker Compose

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8082:8082"
    depends_on:
      - postgres
    environment:
      - CONFIG_PATH=/app/config/local.yaml
    volumes:
      - ./config:/app/config

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: events_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d

volumes:
  postgres_data:
```

---

## API Endpoints

### Аутентификация (публичные)

| Метод | URL             | Описание                    |
|-------|-----------------|----------------------------|
| POST  | /auth/register  | Регистрация пользователя   |
| POST  | /auth/login     | Вход и получение JWT токена|

### Мероприятия (требуют JWT)

| Метод | URL           | Описание                         |
|-------|---------------|----------------------------------|
| POST  | /events       | Создать мероприятие              |
| GET   | /events       | Получить все мероприятия         |
| GET   | /events/{id}  | Получить мероприятие по ID       |

### Примеры запросов

**Регистрация:**
```bash
curl -X POST http://localhost:8082/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","name":"John Doe","password":"password123"}'
```

**Логин:**
```bash
curl -X POST http://localhost:8082/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

**Создание мероприятия:**
```bash
curl -X POST http://localhost:8082/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Go Meetup",
    "description": "Встреча Go разработчиков",
    "location": "Москва, ул. Примерная 1",
    "start_time": "2024-12-20T18:00:00Z",
    "end_time": "2024-12-20T21:00:00Z",
    "max_slots": 50
  }'
```

**Получить все мероприятия:**
```bash
curl -X GET "http://localhost:8082/events?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Получить мероприятие по ID:**
```bash
curl -X GET http://localhost:8082/events/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
