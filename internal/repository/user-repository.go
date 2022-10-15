package repository

import (
    "context"
    "fmt"
    "strings"

    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
    "github.com/jackc/pgx/v5"
)

type UserRepository struct {
    db *pgx.Conn
}

func newUserRepository(db *pgx.Conn) *UserRepository {
    return &UserRepository{
        db: db,
    }
}

func (r *UserRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (User, error) {
    user := User{}

    err := r.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT u.id, u.username, u.name, u.birthday, u.email, u.phone, u.created_at
        FROM %s AS r
        JOIN %s AS u
        ON u.id = r.user_id
        WHERE r.token = $1
    `, refreshTable, usersTable), refreshToken).Scan(
        &user.Id, &user.Username, &user.Name, &user.Birthday,
        &user.Email, &user.Phone, &user.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) {
        return User{}, errors.ErrUsernameNotFound
    }
    if err != nil {
        return User{}, err
    }

    return user, nil
}

func (r *UserRepository) GetByCredentials(ctx context.Context, username string, passHash string) (User, error) {
    user := User{}

    err := r.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT u.id, u.username, u.name, u.birthday, u.email, u.phone, u.created_at
        FROM %s AS a
        JOIN %s AS u
        ON u.id = a.user_id
        WHERE a.password = $1 AND u.username = $2
    `, authTable, usersTable), passHash, username).Scan(
        &user.Id, &user.Username, &user.Name, &user.Birthday,
        &user.Email, &user.Phone, &user.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) {
        return User{}, errors.ErrWrongPassOrUsername
    }
    if err != nil {
        return User{}, err
    }

    return user, nil
}
