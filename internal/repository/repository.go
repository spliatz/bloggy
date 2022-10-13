package repository

import (
    "github.com/jackc/pgx/v5"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
)

const (
    usersTable   = "users"
    authTable    = "auth"
    refreshTable = "refresh"
    postsTable   = "posts"
)

type Repository struct {
    User
    Auth
}

type User interface {
    Create(input *structs.UserCreateInput) (*structs.User, error)
}

type Auth interface {
    Registration(input *structs.UserCreateInput) (int, error)
    AddRefreshToken(input *structs.AuthInput, token string) error
}

func NewRepository(db *pgx.Conn) *Repository {
    return &Repository{
        User: newUserRepository(db),
        Auth: newAuthRepository(db),
    }
}
