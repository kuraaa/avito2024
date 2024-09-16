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