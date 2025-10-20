// Интерфейс и реализация In-memory Cache
// Заменил map на LRU-cache с лимитом и добавил синхронизацию.
// Конструктор возвращает ошибку вместо log.Fatal чтобы не завершать процесс в библиотеке.
package repositories

import (
	"context"
	"fmt"
	"log"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/tvbondar/go-server/internal/entities"
)

type CacheOrderRepository struct {
	cache *lru.Cache[string, entities.Order]
	mu    sync.RWMutex
}

func NewCacheOrderRepository() (*CacheOrderRepository, error) {
	cache, err := lru.New[string, entities.Order](1000)
	if err != nil {
		return nil, err
	}
	return &CacheOrderRepository{cache: cache}, nil
}

func (r *CacheOrderRepository) SaveOrder(ctx context.Context, order entities.Order) error {
	r.mu.Lock()
	r.cache.Add(order.OrderUID, order)
	r.mu.Unlock()
	return nil
}

func (r *CacheOrderRepository) GetOrderByID(ctx context.Context, id string) (entities.Order, error) {
	r.mu.RLock()
	order, ok := r.cache.Get(id)
	r.mu.RUnlock()
	if !ok {
		return entities.Order{}, fmt.Errorf("order not found")
	}
	return order, nil
}

func (r *CacheOrderRepository) GetAllOrders(ctx context.Context) ([]entities.Order, error) {
	var orders []entities.Order
	r.mu.RLock()
	keys := r.cache.Keys()
	r.mu.RUnlock()
	for _, key := range keys {
		r.mu.RLock()
		order, ok := r.cache.Get(key)
		r.mu.RUnlock()
		if !ok {
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// LoadFromDB загружает часть данных в кэш. При большом количестве заказов стоит
// загружать порциями или не грузить всё сразу. Используем передаваемый контекст.
func (r *CacheOrderRepository) LoadFromDB(ctx context.Context, dbRepo OrderRepository) error {
	orders, err := dbRepo.GetAllOrders(ctx)
	if err != nil {
		return err
	}
	for _, order := range orders {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			r.mu.Lock()
			r.cache.Add(order.OrderUID, order)
			r.mu.Unlock()
		}
	}
	log.Println("Cache loaded with", len(orders), "orders")
	return nil
}
