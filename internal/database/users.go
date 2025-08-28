package database

import (
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt" binding:"required, datetime=2006-01-02"`
	UpdatedAt time.Time `json:"updatedAt" binding:"required, datetime=2006-01-02"`
}
