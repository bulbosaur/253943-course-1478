package repository

import (
	"context"
	"errors"
	"strconv"

	pb "lyceum/pkg/api/test"
)

var (
	ErrorOrderID       = errors.New("invalid order id format")
	ErrorOrderNotFound = errors.New("order not found")
)

func (s *PostgresOrderRepository) CreateOrder(ctx context.Context, item string, quantity int32) (string, error) {
	var newID int32

	sql := `INSERT INTO public.orders (item, quantity) VALUES ($1, $2) RETURNING public.orders.id`

	err := s.pool.QueryRow(ctx, sql, item, quantity).Scan(&newID)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(newID)), err
}

func (s *PostgresOrderRepository) GetOrder(ctx context.Context, strID string) (*pb.Order, error) {
	var (
		order pb.Order
		id    int32
	)

	intID, err := strconv.ParseInt(strID, 10, 32)
	if err != nil {
		return nil, ErrorOrderID
	}

	sql := `SELECT id, item, quantity FROM orders WHERE id = $1`

	err = s.pool.QueryRow(ctx, sql, intID).Scan(&id, &order.Item, &order.Quantity)
	if err != nil {
		return nil, err
	}

	order.Id = strconv.Itoa(int(id))
	return &order, nil
}

func (s *PostgresOrderRepository) DeleteOrder(ctx context.Context, strID string) (bool, error) {
	intID, err := strconv.ParseInt(strID, 10, 32)
	if err != nil {
		return false, ErrorOrderID
	}

	sql := `DELETE FROM orders WHERE id = $1`

	res, err := s.pool.Exec(ctx, sql, intID)
	if err != nil {
		return false, err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return false, ErrorOrderNotFound
	}

	return true, nil
}

func (s *PostgresOrderRepository) UpdateOrder(
	ctx context.Context,
	strID, item string,
	quantity int32,
) (*pb.Order, error) {
	var (
		order pb.Order
	)

	intID, err := strconv.ParseInt(strID, 10, 32)
	if err != nil {
		return nil, ErrorOrderID
	}

	sql := `UPDATE orders SET item = $1, quantity = $2 WHERE id = $3`

	res, err := s.pool.Exec(ctx, sql, item, quantity, intID)
	if err != nil {
		return nil, err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, ErrorOrderNotFound
	}

	sql = `SELECT item, quantity FROM orders WHERE id = $1`
	err = s.pool.QueryRow(ctx, sql, intID).Scan(&order.Item, &order.Quantity)
	if err != nil {
		return nil, err
	}

	order.Id = strconv.Itoa(int(intID))

	return &order, nil
}

func (s *PostgresOrderRepository) ListOrders(ctx context.Context) ([]*pb.Order, error) {
	var (
		id     int32
		orders []*pb.Order
	)

	sql := `SELECT id, item, quantity FROM orders`

	rows, err := s.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order pb.Order
		err = rows.Scan(&id, &order.Item, &order.Quantity)
		if err != nil {
			return nil, err
		}

		order.Id = strconv.Itoa(int(id))
		orders = append(orders, &order)
	}

	return orders, nil
}
