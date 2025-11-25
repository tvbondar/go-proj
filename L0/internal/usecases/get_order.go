// Логика получения заказа (из кэша/БД)
// log.Printf для ошибок
package usecases

import (
	"context"
	"log"

	"github.com/tvbondar/go-server/internal/entities"
	"github.com/tvbondar/go-server/internal/repositories"
)

type GetOrderUseCase struct {
	cacheRepo repositories.OrderRepository
	dbRepo    repositories.OrderRepository
}

func NewGetOrderUseCase(cacheRepo, dbRepo repositories.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{cacheRepo: cacheRepo, dbRepo: dbRepo}
}

func (u *GetOrderUseCase) Execute(ctx context.Context, id string) (entities.Order, error) {
	order, err := u.cacheRepo.GetOrderByID(ctx, id)
	if err == nil {
		return order, nil
	}
	log.Printf("Cache miss for order %s: %v", id, err) // Логируем miss для мониторинга

	order, err = u.dbRepo.GetOrderByID(ctx, id)
	if err != nil {
		log.Printf("DB error for order %s: %v", id, err)
		return entities.Order{}, err
	}
	if err := u.cacheRepo.SaveOrder(ctx, order); err != nil {
		log.Printf("Failed to cache order %s: %v", id, err) // Не фатально, но логируем
	}
	return order, nil
}
