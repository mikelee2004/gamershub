FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/server/main.go

FROM alpine:latest

COPY --from=builder /app/bin/server /app/server

# Рабочая директория
WORKDIR /app

EXPOSE 4040


CMD ["/app/server"]