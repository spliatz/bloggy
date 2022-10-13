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
}

type Auth interface {
    SignUp(input *structs.SignUpRequest) (id int, err error)
    AddRefreshToken(input *structs.AuthInput, token string) error
}

func NewRepository(db *pgx.Conn) *Repository {
    return &Repository{
        User: newUserRepository(db),
        Auth: newAuthRepository(db),
    }
}
