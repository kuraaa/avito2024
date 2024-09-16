# Используем базовый образ Go
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Устанавливаем зависимости
RUN go mod download

# Собираем приложение
RUN go build -o main ./cmd/main.go

# Используем образ для выполнения миграций
FROM migrate/migrate:latest AS migrator

# Копируем миграции
COPY --from=builder /app/migrations /migrations

# Копируем скрипт ожидания подключения к базе данных
COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh

# Устанавливаем драйвер PostgreSQL
RUN apk add --no-cache postgresql-client

# Используем образ для запуска приложения
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение
COPY --from=builder /app/main .

# Копируем файл .env
COPY .env .

# Команда для запуска приложения
CMD ["./main"]