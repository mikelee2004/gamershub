# Билд стадия
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем только модули сначала для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /gamershub ./cmd/server/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем только бинарник
COPY --from=builder /gamershub /app/

# Создаем пустой .env если его нет (значения будут из docker-compose)
RUN touch .env

EXPOSE 8080

CMD ["./gamershub"]