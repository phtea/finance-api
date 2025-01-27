# Используем golang image для сборки приложения
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum для управления зависимостями
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod tidy

# Копируем весь код приложения
COPY . .

# Строим приложение
RUN go build -o main cmd/main/main.go

# Используем легкий образ Alpine для финального контейнера
FROM alpine:latest

# Устанавливаем необходимые для работы библиотеки
RUN apk --no-cache add ca-certificates

# Копируем скомпилированное приложение из стадии сборки
COPY --from=builder /app/main /main

# Делаем приложение исполнимым
RUN chmod +x /main

# Запускаем приложение
CMD ["/main"]
