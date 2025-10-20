# builder
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o server ./cmd

# final
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8081
ENTRYPOINT ["./server"]