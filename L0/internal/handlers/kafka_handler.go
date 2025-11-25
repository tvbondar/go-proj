// Логика обработки сообщений из Kafka
// Принимаем cfg извне, возвращаем ошибку в случае фатального сбоя.
// Ctx используется для graceful shutdown.
package handlers

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/tvbondar/go-server/internal/config"
	"github.com/tvbondar/go-server/internal/usecases"
)

func StartKafkaConsumer(ctx context.Context, cfg *config.Config, usecase *usecases.ProcessOrderUseCase) error {

	readerCfg := kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaAddr},
		Topic:    "orders",
		GroupID:  "my-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	}

	var reader *kafka.Reader

	// Попытки подключения к брокеру
	for i := 0; i < 10; i++ {
		// Проверяем доступность брокера через Dial
		conn, err := kafka.Dial("tcp", cfg.KafkaAddr)
		if err == nil {
			_ = conn.Close()
			// создаём reader только если брокер доступен
			reader = kafka.NewReader(readerCfg)
			break
		}
		log.Printf("Failed to connect to Kafka (attempt %d): %v", i+1, err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	if reader == nil {
		return nil // не фатально: если Kafka недоступен, возвращаем и позволяем main продолжить работу
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing kafka reader: %v", err)
		}
	}()
	log.Println("Successfully connected to Kafka topic 'orders'")

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer shutting down")
			return nil
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				// Если ctx отменён, выйдем; иначе залогируем и продолжим
				if ctx.Err() != nil {
					return ctx.Err()
				}
				log.Printf("Error fetching message: %v", err)
				continue
			}
			log.Printf("Received message with key: %s", string(msg.Key))
			if err := usecase.Execute(ctx, msg.Value); err != nil {
				log.Printf("Error processing message: %v", err)
				// Не коммитим при ошибке — позволяем retry
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
