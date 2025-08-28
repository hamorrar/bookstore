package database

import "database/sql"

type Models struct {
	Users  UserModel
	Orders OrderModel
	Books  BookModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:  UserModel{DB: db},
		Orders: OrderModel{DB: db},
		Books:  BookModel{DB: db},
	}
}
