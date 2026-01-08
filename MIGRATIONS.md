# Миграции базы данных (PostgreSQL)

Этот документ содержит SQL скрипты для настройки базы данных PostgreSQL.

## Создание базы данных

```sql
CREATE DATABASE events_db;
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

## Миграция 002: Таблица users (с профилем и балансом)

```sql
-- 002_create_users_table.sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar_url TEXT,
    bio TEXT,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
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

## Миграция 003: Таблица events (мероприятия с расширенными полями)

```sql
-- 003_create_events_table.sql
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL DEFAULT 'other',
    image_url TEXT,
    venue VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    capacity INTEGER NOT NULL DEFAULT 1,
    available_tickets INTEGER NOT NULL DEFAULT 1,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_events_creator_id ON events(creator_id);
CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);
CREATE INDEX IF NOT EXISTS idx_events_category ON events(category);
CREATE INDEX IF NOT EXISTS idx_events_price ON events(price);

-- Полнотекстовый поиск (для русского и английского языков)
CREATE INDEX IF NOT EXISTS idx_events_search ON events USING GIN (
    to_tsvector('russian', title || ' ' || COALESCE(description, '') || ' ' || venue || ' ' || COALESCE(address, ''))
);

-- Триггер для автоматического обновления updated_at
CREATE TRIGGER update_events_updated_at 
    BEFORE UPDATE ON events
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### Категории мероприятий

| Значение | Описание |
|----------|----------|
| concert | Концерты |
| sport | Спортивные мероприятия |
| theater | Театр |
| exhibition | Выставки |
| festival | Фестивали |
| other | Другое |

---

## Миграция 004: Таблица bookings (бронирования)

```sql
-- 004_create_bookings_table.sql
CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    total_price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed',
    booking_code VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Один пользователь может забронировать одно мероприятие только один раз
    UNIQUE(user_id, event_id)
);

CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_event_id ON bookings(event_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_booking_code ON bookings(booking_code);
```

### Статусы бронирования

| Значение | Описание |
|----------|----------|
| confirmed | Подтверждено |
| cancelled | Отменено |
| used | Использовано |

---

## Применение всех миграций

```bash
# Все миграции одной командой
psql -h localhost -U postgres -d events_db << 'EOF'

-- Миграция 001
CREATE TABLE IF NOT EXISTS url (
    id BIGSERIAL PRIMARY KEY,
    alias VARCHAR(255) NOT NULL UNIQUE,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_url_alias ON url(alias);

-- Миграция 002
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar_url TEXT,
    bio TEXT,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Миграция 003
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL DEFAULT 'other',
    image_url TEXT,
    venue VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    capacity INTEGER NOT NULL DEFAULT 1,
    available_tickets INTEGER NOT NULL DEFAULT 1,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_events_creator_id ON events(creator_id);
CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);
CREATE INDEX IF NOT EXISTS idx_events_category ON events(category);
CREATE INDEX IF NOT EXISTS idx_events_price ON events(price);
CREATE INDEX IF NOT EXISTS idx_events_search ON events USING GIN (
    to_tsvector('russian', title || ' ' || COALESCE(description, '') || ' ' || venue || ' ' || COALESCE(address, ''))
);

CREATE TRIGGER update_events_updated_at 
    BEFORE UPDATE ON events
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Миграция 004
CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    total_price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed',
    booking_code VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, event_id)
);

CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_event_id ON bookings(event_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_booking_code ON bookings(booking_code);

EOF
```

---

## Docker Compose

```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8082:8082"
    depends_on:
      postgres:
        condition: service_healthy
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
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

---

## API Endpoints

### Аутентификация (публичные)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | /auth/register | Регистрация пользователя |
| POST | /auth/login | Вход и получение JWT токена |

### Профиль (требует JWT)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | /profile | Получить профиль |
| PUT | /profile | Обновить профиль |
| POST | /profile/balance | Пополнить баланс |

### Мероприятия (требует JWT)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | /events | Создать мероприятие |
| GET | /events | Получить все мероприятия |
| GET | /events/{id} | Получить мероприятие по ID |
| POST | /events/{id}/book | Забронировать билет |

### Бронирования (требует JWT)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | /bookings | Мои билеты |
| DELETE | /bookings/{id} | Отменить бронь |

### Поиск (требует JWT)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | /search | Поиск мероприятий |

---

## Примеры запросов

### Регистрация
```bash
curl -X POST http://localhost:8082/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","name":"John Doe","password":"password123"}'
```

### Логин
```bash
curl -X POST http://localhost:8082/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Получить профиль
```bash
curl -X GET http://localhost:8082/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Обновить профиль
```bash
curl -X PUT http://localhost:8082/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name":"John Updated","phone":"+79001234567","bio":"Developer"}'
```

### Пополнить баланс
```bash
curl -X POST http://localhost:8082/profile/balance \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"amount": 5000.00}'
```

### Создание мероприятия
```bash
curl -X POST http://localhost:8082/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Рок-концерт 2024",
    "description": "Лучший рок-концерт года!",
    "category": "concert",
    "venue": "Олимпийский",
    "address": "Москва, Олимпийский проспект, 16",
    "price": 3500.00,
    "capacity": 15000,
    "start_time": "2024-12-20T19:00:00Z",
    "end_time": "2024-12-20T23:00:00Z"
  }'
```

### Забронировать билет
```bash
curl -X POST http://localhost:8082/events/1/book \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"quantity": 2}'
```

### Мои билеты
```bash
curl -X GET http://localhost:8082/bookings \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Отменить бронь
```bash
curl -X DELETE http://localhost:8082/bookings/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Поиск мероприятий
```bash
# Простой поиск
curl -X GET "http://localhost:8082/search?q=концерт" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# С фильтрами
curl -X GET "http://localhost:8082/search?q=москва&category=concert&price_min=1000&price_max=5000&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# По датам
curl -X GET "http://localhost:8082/search?date_from=2024-12-01T00:00:00Z&date_to=2024-12-31T23:59:59Z" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Откат миграций

```sql
-- Откат миграции 004
DROP TABLE IF EXISTS bookings;

-- Откат миграции 003
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP TABLE IF EXISTS events;

-- Откат миграции 002
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS users;

-- Откат миграции 001
DROP TABLE IF EXISTS url;
