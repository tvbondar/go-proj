// Интерфейс Repository (общий для кэша и БД)
// Context позволяет отменять операции (например, при shutdown).
// Пока не меняй реализации — добавь ctx позже, когда реализуешь shutdown в main.go.
package repositories

import (
	"context"

	"github.com/tvbondar/go-server/internal/entities"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order entities.Order) error
	GetOrderByID(ctx context.Context, id string) (entities.Order, error)
	GetAllOrders(ctx context.Context) ([]entities.Order, error) // Для восстановления кэша
}
