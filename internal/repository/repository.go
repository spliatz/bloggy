package repository

import (
    "github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
    "github.com/jackc/pgx/v5"
)

const (
    usersTable = "users"
    authTable  = "auth"
    postsTable = "posts"
)

type Repository struct {
    User
}

type User interface {
    Create(input *structs.UserCreateInput) (*structs.User, error)
}

func NewRepository(db *pgx.Conn) *Repository {
    return &Repository{
        User: newUserRepository(db),
    }
}
