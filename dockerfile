# Строим с использованием golang:1.22-alpine
FROM golang:1.22-alpine as builder

# Устанавливаем рабочую директорию для сборки
WORKDIR /app

# Копируем go.mod и go.sum в рабочую директорию
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем весь проект в контейнер
COPY . .

# Переходим в директорию, где находится main.go
WORKDIR /app/cmd/app

# Собираем исполняемый файл с именем main
RUN go build -o main .

# Второй этап: переносим только исполнимый файл в более легкий образ
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный файл из предыдущего этапа
COPY --from=builder /app/cmd/app/main .

COPY .env .env

# Устанавливаем точку входа
CMD ["./main"]