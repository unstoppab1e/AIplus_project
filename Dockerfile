# Используем многоступенчатую сборку
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем все исходники сразу
COPY . .

# Подтягиваем зависимости (сразу с исходниками)
RUN go mod tidy
RUN go mod download

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем собранное приложение и миграции
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]
