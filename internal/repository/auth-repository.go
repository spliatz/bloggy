package repository

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgtype"

    e "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type AuthRepository struct {
    db *pgx.Conn
}

func newAuthRepository(db *pgx.Conn) *AuthRepository {
    return &AuthRepository{
        db: db,
    }
}

type User struct {
    Id        int         `db:"id"`
    Username  string      `db:"username"`
    Password  string      `db:"password"`
    Name      pgtype.Text `db:"name"`
    Birthday  pgtype.Date `db:"birthday"`
    Email     pgtype.Text `db:"email"`
    Phone     pgtype.Text `db:"phone"`
    CreatedAt time.Time   `db:"created_at"`
}

func (r *AuthRepository) SignUp(ctx context.Context, u User) (User, error) {

    _id := 0
    err := r.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE username = $1
    `, usersTable), u.Username).Scan(&_id)
    if err == nil {
        // Если он нашел пользователя и успешно просканировал, то он существует
        return User{}, e.ErrTakenUsername
    }
    if !errors.Is(err, pgx.ErrNoRows) {
        // Если ошибка не является ошибкой "Не найден пользователь"
        return User{}, err
    }

    tx, err := r.db.Begin(ctx)

    newU := User{}

    // FIXME: Тут нет проверки занятости email и phone
    // Создание пользователя
    err = tx.QueryRow(ctx, fmt.Sprintf(`
        INSERT INTO %s (username, name, email, phone, birthday, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, username, name, email, phone, birthday, created_at
    `, usersTable), u.Username, u.Name, u.Email, u.Phone, u.Birthday, u.CreatedAt).Scan(
        &newU.Id, &newU.Username, &newU.Name, &newU.Email, &newU.Phone, &newU.Birthday, &newU.CreatedAt)
    if err != nil {
        _ = tx.Rollback(ctx)
        return User{}, err
    }

    // Привязка пароля
    err = tx.QueryRow(ctx, fmt.Sprintf(`
        INSERT INTO %s (user_id, password)
        VALUES ($1, $2)
        RETURNING user_id, password
    `, authTable), newU.Id, u.Password).Scan(&newU.Id, &newU.Password)
    if err != nil {
        _ = tx.Rollback(ctx)
        return User{}, err
    }

    err = tx.Commit(ctx)
    if err != nil {
        return User{}, err
    }

    return newU, nil
}

type Session struct {
    RefreshToken string    `db:"token"`
    ExpiresAt    time.Time `db:"expires_at"`
}

func (r *AuthRepository) SetSession(ctx context.Context, userId int, s Session) error {

    _, err := r.db.Exec(ctx, fmt.Sprintf(`
        INSERT INTO %s(user_id, token, expires_at)
        VALUES ($1, $2, $3)
    `, refreshTable), userId, s.RefreshToken, s.ExpiresAt)
    if err != nil {
        return err
    }

    return nil
}

func (r *AuthRepository) CheckRefresh(ctx context.Context, refreshToken string) error {
    var expiresAt pgtype.Date
    err := r.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT expires_at
        FROM %s
        WHERE token = $1
    `, refreshTable), refreshToken).Scan(&expiresAt)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return e.ErrTokenNotFound
        }

        return err
    }

    if expiresAt.Valid {
        if time.Now().After(expiresAt.Time) {
            return e.ErrTokenExpired
        }
    }

    return nil
}

func (r *AuthRepository) DeleteRefresh(ctx context.Context, refreshToken string) error {
    err := r.db.QueryRow(ctx, fmt.Sprintf(`
                DELETE FROM %s
                WHERE token = $1
                RETURNING token
            `, refreshTable), refreshToken).Scan(&refreshToken)
    if err != nil && !errors.Is(err, pgx.ErrNoRows) {
        return err
    }
    if errors.Is(err, pgx.ErrNoRows) {
        return e.ErrTokenNotFound
    }

    return nil
}
