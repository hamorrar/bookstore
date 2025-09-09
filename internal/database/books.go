package database

import (
	"context"
	"database/sql"
	"time"
)

type BookModel struct {
	DB *sql.DB
}

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title" binding:"required,min=3"`
	Author string `json:"author" binding:"required,min=3"`
	Price  int    `json:"price" binding:"required"`
}

func (m *BookModel) CreateBook(book *Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "insert into books (book_title, book_author, book_price) values ($1, $2, $3) returning book_id"

	return m.DB.QueryRowContext(ctx, query, book.Title, book.Author, book.Price).Scan(&book.Id)
}

func (m *BookModel) DeleteBook(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "delete from books where book_id = $1"

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *BookModel) GetBook(id int) (*Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select * from books where book_id = $1"

	var book Book

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&book.Id, &book.Title, &book.Author, &book.Price)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

func (m *BookModel) GetPageOfBooks(limit int, page int) ([]*Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 2
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	query := "select * from books order by book_id limit $1 offset $2"

	rows, err := m.DB.QueryContext(ctx, query, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := []*Book{}

	for rows.Next() {
		var book Book

		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Price)

		if err != nil {
			return nil, err
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (m *BookModel) UpdateBook(book *Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "update books set book_title = $1, book_author = $2, book_price = $3 where book_id = $4"

	_, err := m.DB.ExecContext(ctx, query, book.Title, book.Author, book.Price, book.Id)

	if err != nil {
		return err
	}
	return nil
}
