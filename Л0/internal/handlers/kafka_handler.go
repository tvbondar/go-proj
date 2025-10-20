// Логика обработки сообщений из Kafka
// Вынес конфиг в LoadConfig. Добавил ctx для shutdown. Заменил fmt на log.
// Конфиг из config.go — гибко. Ctx для shutdown (вызови go handlers.StartKafkaConsumer(ctx, ...) в main.go).
// При ошибке обработки не коммитим — retry от Kafka. Это предотвращает потерю данных.
package handlers

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/tvbondar/go-server/internal/config"
	"github.com/tvbondar/go-server/internal/usecases"
)

func StartKafkaConsumer(ctx context.Context, usecase *usecases.ProcessOrderUseCase) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config for Kafka:", err)
	}

	readerCfg := kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaAddr},
		Topic:    "orders",
		GroupID:  "my-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	}

	var reader *kafka.Reader

	for i := 0; i < 10; i++ {
		reader = kafka.NewReader(readerCfg)
		var conn *kafka.Conn
		conn, err = kafka.Dial("tcp", cfg.KafkaAddr)
		if err == nil {
			conn.Close()
			break
		}
		log.Printf("Failed to connect to Kafka (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Printf("Failed to connect to Kafka after retries: %v", err)
		return
	}
	defer reader.Close()
	log.Println("Successfully connected to Kafka topic 'orders'")

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer shutting down")
			return
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error fetching message: %v", err)
				continue
			}
			log.Printf("Received message with key: %s", string(msg.Key))
			if err := usecase.Execute(ctx, msg.Value); err != nil {
				log.Printf("Error processing message: %v", err)
				// Можно не коммитить при ошибке, чтобы retry
				continue
			}
			if err := reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Error committing message: %v", err)
			} else {
				log.Printf("Message processed and committed")
			}
		}
	}
}
