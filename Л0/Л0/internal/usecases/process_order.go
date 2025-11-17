// Логика обработки нового заказа (сохранение)
// ручную проверку на validator — это автоматически проверит все поля с тегами (добавь теги в order.go).
// Если валидация fails, возвращаем err (для логирования в Kafka). Логируем с log.
package usecases

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/tvbondar/go-server/internal/entities"
	"github.com/tvbondar/go-server/internal/repositories"
)

type ProcessOrderUseCase struct {
	dbRepo    repositories.OrderRepository
	cacheRepo repositories.OrderRepository
}

func NewProcessOrderUseCase(dbRepo, cacheRepo repositories.OrderRepository) *ProcessOrderUseCase {
	return &ProcessOrderUseCase{dbRepo: dbRepo, cacheRepo: cacheRepo}
}

func (u *ProcessOrderUseCase) Execute(ctx context.Context, rawMessage []byte) error {
	var order entities.Order
	if err := json.Unmarshal(rawMessage, &order); err != nil {
		log.Printf("Invalid JSON message: %v", err)
		return err // Возвращаем err, чтобы не коммитить в Kafka
	}

	validate := validator.New()
	if err := validate.Struct(order); err != nil {
		log.Printf("Validation failed for order %s: %v", order.OrderUID, err)
		return err // Игнорируем невалидные, но логируем и не сохраняем
	}

	if err := u.dbRepo.SaveOrder(ctx, order); err != nil {
		log.Printf("DB save error for order %s: %v", order.OrderUID, err)
		return err
	}
	if err := u.cacheRepo.SaveOrder(ctx, order); err != nil {
		log.Printf("Cache save error for order %s: %v", order.OrderUID, err)
		// Не фатально, продолжаем
	}
	return nil
}
