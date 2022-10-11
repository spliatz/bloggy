package repository

import (
	"context"
	"github.com/jackc/pgx"
)

type UserRepository struct {
	db *pgx.Conn
}

func newUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (s *UserRepository) Create(username, name, surname, email, password string) error {
	return s.db.Ping(context.Background())
}
