package database

import (
	"database/sql"
)

type OrderModel struct {
	DB *sql.DB
}

type Order struct {
	Id         int    `json:"id"`
	Userid     int    `json:"userId" binding:"required"`
	Status     string `json:"status"`
	TotalPrice int    `json:"totalPrice"`
}
