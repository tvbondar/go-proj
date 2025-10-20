// Работа с PostgreSQL
// Перевёл операции на использование контекста (BeginTx, ExecContext, QueryRowContext)
package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tvbondar/go-server/internal/entities"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) SaveOrder(ctx context.Context, order entities.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	deliveryJSON, err := json.Marshal(order.Delivery)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery: %w", err)
	}
	paymentJSON, err := json.Marshal(order.Payment)
	if err != nil {
		return fmt.Errorf("failed to marshal payment: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery, payment)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
		deliveryJSON, paymentJSON)
	if err != nil {
		return fmt.Errorf("failed to insert order %s: %w", order.OrderUID, err)
	}

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("failed to insert item for order %s: %w", order.OrderUID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx for order %s: %w", order.OrderUID, err)
	}
	return nil
}

func (r *PostgresOrderRepository) GetOrderByID(ctx context.Context, id string) (entities.Order, error) {
	var order entities.Order
	var deliveryJSON, paymentJSON []byte

	err := r.db.QueryRowContext(ctx, `
        SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery, payment
        FROM orders WHERE order_uid = $1`,
		id).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
		&deliveryJSON, &paymentJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Order{}, fmt.Errorf("order not found")
		}
		return entities.Order{}, fmt.Errorf("failed to query order %s: %w", id, err)
	}

	if err := json.Unmarshal(deliveryJSON, &order.Delivery); err != nil {
		return entities.Order{}, fmt.Errorf("failed to unmarshal delivery: %w", err)
	}
	if err := json.Unmarshal(paymentJSON, &order.Payment); err != nil {
		return entities.Order{}, fmt.Errorf("failed to unmarshal payment: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
        SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
        FROM items WHERE order_uid = $1`, id)
	if err != nil {
		return entities.Order{}, fmt.Errorf("failed to query items for order %s: %w", id, err)
	}
	defer rows.Close()
	for rows.Next() {
		var item entities.Item
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return entities.Order{}, fmt.Errorf("failed to scan item for order %s: %w", id, err)
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *PostgresOrderRepository) GetAllOrders(ctx context.Context) ([]entities.Order, error) {
	var orders []entities.Order

	// Оптимизированный запрос: все orders + items в одном проходе
	rows, err := r.db.QueryContext(ctx, `
		SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, 
		       o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard, o.delivery, o.payment,
		       i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
		FROM orders o
		LEFT JOIN items i ON o.order_uid = i.order_uid
		ORDER BY o.order_uid`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all orders: %w", err)
	}
	defer rows.Close()

	currentUID := ""
	var currentOrder entities.Order
	for rows.Next() {
		var order entities.Order
		var deliveryJSON, paymentJSON []byte
		var item entities.Item
		var itemChrtID sql.NullInt64

		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
			&deliveryJSON, &paymentJSON,
			&itemChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			log.Printf("Scan error in GetAllOrders: %v", err)
			continue
		}

		if order.OrderUID != currentUID {
			if currentUID != "" {
				orders = append(orders, currentOrder)
			}
			currentUID = order.OrderUID
			currentOrder = order
			if err := json.Unmarshal(deliveryJSON, &currentOrder.Delivery); err != nil {
				log.Printf("Unmarshal delivery error: %v", err)
				continue
			}
			if err := json.Unmarshal(paymentJSON, &currentOrder.Payment); err != nil {
				log.Printf("Unmarshal payment error: %v", err)
				continue
			}
			currentOrder.Items = []entities.Item{}
		}
		if itemChrtID.Valid {
			item.ChrtID = int(itemChrtID.Int64)
			currentOrder.Items = append(currentOrder.Items, item)
		}
	}
	if currentUID != "" {
		orders = append(orders, currentOrder)
	}

	return orders, nil
}
