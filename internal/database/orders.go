package database

import (
	"context"
	"database/sql"
	"time"
)

type OrderModel struct {
	DB *sql.DB
}

type Order struct {
	Id          int    `json:"id"`
	User_Id     int    `json:"user_id" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Total_Price int    `json:"total_price" binding:"required"`
}

func (m *OrderModel) CreateOrder(order *Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "insert into orders (order_user_id, order_status, order_total_price) values ($1, $2, $3) returning order_id"

	err := m.DB.QueryRowContext(ctx, query, order.User_Id, order.Status, order.Total_Price).Scan(&order.Id)
	if err != nil {
		return err
	}
	return nil
}

func (m *OrderModel) DeleteOrder(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "delete from orders where order_id = $1"

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *OrderModel) GetOrder(id int) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select * from orders where order_id = $1"

	var order Order

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&order.Id, &order.User_Id, &order.Status, &order.Total_Price)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (m *OrderModel) GetPageOfOrders(limit int, page int) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 2
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	query := "select * from orders order by order_id limit $1 offset $2"

	rows, err := m.DB.QueryContext(ctx, query, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []*Order{}

	for rows.Next() {
		var order Order

		err := rows.Scan(&order.Id, &order.User_Id, &order.Status, &order.Total_Price)

		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (m *OrderModel) GetAllOrders() ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select * from orders"

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []*Order{}

	for rows.Next() {
		var order Order

		err := rows.Scan(&order.Id, &order.User_Id, &order.Status, &order.Total_Price)

		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (m *OrderModel) UpdateOrder(order *Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE orders SET order_user_id = $1, order_status = $2, order_total_price = $3 WHERE order_id = $4"

	_, err := m.DB.ExecContext(ctx, query, order.User_Id, order.Status, order.Total_Price, order.Id)

	if err != nil {
		return err
	}

	return nil
}
