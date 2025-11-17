# builder
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Зависимости сначала
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники и собираем бинарь
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o server ./cmd

# final
FROM alpine:3.18
WORKDIR /app

# Копируем бинарь из билдера
COPY --from=builder /app/server ./

# Копируем миграции и статику в финальный образ, чтобы goose и http.FileServer работали внутри контейнера
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/web ./web
COPY --from=builder /app/config.yaml ./config.yaml

# Создадим непривилегированного пользователя 
RUN addgroup -S app && adduser -S app -G app || true
USER app

EXPOSE 8081

ENTRYPOINT ["./server"]