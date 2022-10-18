package repository

import (
    "context"
    "errors"
    "fmt"
    "strings"

    e "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
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

func (r *UserRepository) EditById(ctx context.Context, id int, req map[string]string) (User, error) {
    // Проверяем доступность username
    _id := 0
    err := r.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE username = $1
    `, usersTable), req["username"]).Scan(&_id)
    if err == nil {
        // Если он нашел пользователя и успешно просканировал, то он существует
        return User{}, e.ErrTakenUsername
    }
    if !errors.Is(err, pgx.ErrNoRows) {
        // Если ошибка не является ошибкой "Не найден пользователь"
        return User{}, err
    }

    // Проверяем доступность email
    err = r.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE email = $1
    `, usersTable), req["email"]).Scan(&_id)
    if err == nil {
        // Если он нашел пользователя и успешно просканировал, то он существует
        return User{}, e.ErrTakenEmail
    }
    if !errors.Is(err, pgx.ErrNoRows) {
        // Если ошибка не является ошибкой "Не найден пользователь"
        return User{}, err
    }

    // Проверяем доступность phone
    err = r.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE phone = $1
    `, usersTable), req["phone"]).Scan(&_id)
    if err == nil {
        // Если он нашел пользователя и успешно просканировал, то он существует
        return User{}, e.ErrTakenPhone
    }
    if !errors.Is(err, pgx.ErrNoRows) {
        // Если ошибка не является ошибкой "Не найден пользователь"
        return User{}, err
    }

    // Создаем параметры
    var updateQuery []string
    for key, value := range req {
        updateQuery = append(updateQuery, fmt.Sprintf(`%s = '%s'`, key, value))
    }

    user := User{}
    err = r.db.QueryRow(ctx, fmt.Sprintf(`
        UPDATE %s 
        SET %s
        WHERE id = $1
        RETURNING username, name, phone, email, birthday
    `, usersTable, strings.Join(updateQuery, ", ")), id).Scan(
        &user.Username, &user.Name, &user.Phone, &user.Email, &user.Birthday)
    if err != nil {
        return User{}, err
    }

    return user, nil
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
        return User{}, e.ErrUsernameNotFound
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
        return User{}, e.ErrWrongCredentials
    }
    if err != nil {
        return User{}, err
    }

    return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (User, error) {
    user := User{}

    err := r.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id, username, name, birthday, email, phone, created_at
        FROM %s
        WHERE username = $1
    `, usersTable), username).Scan(
        &user.Id, &user.Username, &user.Name, &user.Birthday,
        &user.Email, &user.Phone, &user.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) {
        return User{}, e.ErrUsernameNotFound
    }
    if err != nil {
        return User{}, err
    }

    return user, nil
}
