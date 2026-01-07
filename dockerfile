# Этап сборки (Builder)
FROM golang:alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod ./
# Если есть go.sum, раскомментируйте строку ниже
# COPY go.sum ./

RUN go mod download

# Копируем исходный код
COPY . .

# --- ВАЖНОЕ ИСПРАВЛЕНИЕ ---
# Так как main.go лежит в корне, мы собираем текущую директорию (.)
RUN go build -o main .

# Финальный этап (Production Image)
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/main .

# Копируем конфиги (на случай проблем с volume)
COPY --from=builder /app/config ./config

# Запускаем приложение
CMD ["./main"]