package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"-"`
	Role     string `json:"role" binding:"required"`
}

func (m *UserModel) CreateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "insert into users (user_email, user_password, user_role) values ($1, $2, $3) returning user_id"

	// TODO: hash the password

	err := m.DB.QueryRowContext(ctx, query, user.Email, user.Password, user.Role).Scan(&user.Id)

	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "delete from users where user_id = $1"

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) GetUser(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select * from users where user_id = $1"

	var user User

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Email, &user.Password, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) GetAllUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select * from users"

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Role)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UserModel) UpdateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE users SET user_email = $1, user_password = $2, user_role = $3 where user_id = $4"

	_, err := m.DB.ExecContext(ctx, query, user.Email, user.Password, user.Role, user.Id)

	if err != nil {
		return err
	}

	return nil
}
