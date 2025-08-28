package database

import (
	"database/sql"
	"time"
)

type OrderModel struct {
	DB *sql.DB
}

type Order struct {
	Id         int       `json:"id"`
	Userid     int       `json:"userId" binding:"required"`
	Status     string    `json:"status"`
	TotalPrice int       `json:"totalPrice"`
	CreatedAt  time.Time `json:"createdAt" binding:"required, datetime=2006-01-02"`
	UpdatedAt  time.Time `json:"updatedAt" binding:"required, datetime=2006-01-02"`
}
