package database

import (
	"database/sql"
	"time"
)

type BookModel struct {
	DB *sql.DB
}

type Book struct {
	Id        int       `json:"id"`
	Title     string    `json:"title" binding:"required,min=3"`
	Author    string    `json:"author" binding:"required,min=3"`
	Isbn      int       `json:"isbn" binding:"required,min=3"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"createdAt" binding:"required, datetime=2006-01-02"`
	UpdatedAt time.Time `json:"updatedAt" binding:"required, datetime=2006-01-02"`
}
