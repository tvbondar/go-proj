// Интерфейс и реализация In-memory Cache
// Заменил map на LRU-cache с лимитом (например, 1000 элементов) — старые вытесняются. Сохранил интерфейс
// LRU автоматически удаляет старые записи при переполнении (инвалидация). Это предотвращает утечки памяти
package repositories

import (
	"context"
	"fmt"
	"log"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/tvbondar/go-server/internal/entities"
)

type CacheOrderRepository struct {
	cache *lru.Cache[string, entities.Order]
}

func NewCacheOrderRepository() *CacheOrderRepository {
	cache, err := lru.New[string, entities.Order](1000)
	if err != nil {
		log.Fatal("Failed to create LRU cache:", err)
	}
	return &CacheOrderRepository{cache: cache}
}

func (r *CacheOrderRepository) SaveOrder(ctx context.Context, order entities.Order) error {
	r.cache.Add(order.OrderUID, order)
	return nil
}

func (r *CacheOrderRepository) GetOrderByID(ctx context.Context, id string) (entities.Order, error) {
	order, ok := r.cache.Get(id)
	if !ok {
		return entities.Order{}, fmt.Errorf("order not found")
	}
	return order, nil
}

func (r *CacheOrderRepository) GetAllOrders(ctx context.Context) ([]entities.Order, error) {
	var orders []entities.Order
	keys := r.cache.Keys()
	for _, key := range keys {
		order, _ := r.cache.Get(key)
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *CacheOrderRepository) LoadFromDB(dbRepo OrderRepository) error {
	orders, err := dbRepo.GetAllOrders(context.Background()) // Передаем контекст
	if err != nil {
		return err
	}
	for _, order := range orders {
		r.cache.Add(order.OrderUID, order)
	}
	fmt.Println("Cache loaded with", len(orders), "orders")
	return nil
}
