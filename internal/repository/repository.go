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
    GetByUsername(сtx context.Context, username string) (User, error)
    GetByRefreshToken(сtx context.Context, refreshToken string) (User, error)
    GetByCredentials(ctx context.Context, username string, passHash string) (User, error)
    EditById(сtx context.Context, id int, req map[string]string) (User, error)
}

type Auth interface {
    SignUp(ctx context.Context, u User) (User, error)
    SetSession(ctx context.Context, userId int, s Session) error
    CheckRefresh(ctx context.Context, refresh string) error
    DeleteRefresh(ctx context.Context, refreshToken string) error
}

type PostRepo interface {
    Create(req Post) (int, error)
    GetOne(id int) (Post, error)
    GetAllUserPosts(username string) ([]Post, error)
    Delete(id int) error
}

func NewRepository(db *pgx.Conn) *Repository {
    return &Repository{
        UserRepo: newUserRepository(db),
        Auth:     newAuthRepository(db),
        PostRepo: newPostRepository(db),
    }
}
