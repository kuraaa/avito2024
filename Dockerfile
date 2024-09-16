# Шаг сборки Go приложения
FROM golang:1.20 AS build

WORKDIR /app

# Копируйте файлы go.mod и go.sum из директории backend
COPY backend/go.mod backend/go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируйте остальные файлы проекта
COPY backend/ .

# Соберите приложение Go
RUN go build -o myapp .

# Шаг выполнения: используйте тот же образ golang для совместимости библиотек
FROM golang:1.20

# Скопируйте собранное приложение из этапа сборки
COPY --from=build /app/myapp /myapp

# Укажите порт, который будет прослушиваться
EXPOSE 8080

# Установите команду запуска контейнера
CMD ["/myapp"]
