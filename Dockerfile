# Этап сборки
FROM golang:1.24-alpine AS builder

# Установка зависимостей для сборки
RUN apk add --no-cache git make

# Установка рабочей директории
WORKDIR /app

# Копирование go mod файлов
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Генерация GraphQL кода
RUN go run github.com/99designs/gqlgen generate

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/server ./cmd/server

# Финальный этап - минимальный образ
FROM alpine:latest

# Установка CA сертификатов для HTTPS запросов
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копирование бинарного файла из builder
COPY --from=builder /app/bin/server .

# Копирование миграций (если нужны)
COPY --from=builder /app/migrations ./migrations

# Создание пользователя для запуска приложения
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

# Открытие порта
EXPOSE 8080

# Команда запуска
CMD ["./server"]

