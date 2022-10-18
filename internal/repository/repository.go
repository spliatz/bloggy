package repository

import (
    "context"

    "github.com/jackc/pgx/v5"
)

const (
    usersTable   = "users"
    authTable    = "auth"
    refreshTable = "refresh"
    postsTable   = "posts"
)

type Repository struct {
    UserRepo
    Auth
    PostRepo
}

type UserRepo interface {
    GetByUsername(ctx context.Context, username string) (User, error)
    GetById(ctx context.Context, id int) (User, error)
    GetByRefreshToken(ctx context.Context, refreshToken string) (User, error)
    GetByCredentials(ctx context.Context, username string, passHash string) (User, error)
    EditById(ctx context.Context, id int, req map[string]string) (User, error)
}

type Auth interface {
    SignUp(ctx context.Context, u User) (User, error)
    SetSession(ctx context.Context, userId int, s Session) error
    CheckRefresh(ctx context.Context, refreshToken string) error
    DeleteRefresh(ctx context.Context, refreshToken string) error
}

type PostRepo interface {
    Create(ctx context.Context, p Post) (id int, err error)
    GetById(ctx context.Context, id int) (Post, error)
    GetAllByUsername(ctx context.Context, username string) ([]Post, error)
    DeleteById(ctx context.Context, id int) error
    IsAuthor(ctx context.Context, postId int, authorId int) (bool, error)
}

func NewRepository(db *pgx.Conn) *Repository {
    return &Repository{
        UserRepo: newUserRepository(db),
        Auth:     newAuthRepository(db),
        PostRepo: newPostRepository(db),
    }
}
