package repository

import (
	"github.com/jackc/pgx"
)

type Repository struct {
	User
}

type User interface {
	Create(username, name, surname, email, password string) error
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		User: newUserRepository(db),
	}
}
