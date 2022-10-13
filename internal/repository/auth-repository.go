package repository

import (
    "context"
    "errors"
    "fmt"

    "github.com/jackc/pgx/v5"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
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

func (r *AuthRepository) SignUp(req *structs.SignUpRequest) (id int, err error) {

    err = r.db.QueryRow(context.Background(), fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE username = $1
    `, usersTable), req.Username).Scan(&id)
    if err == nil {
        // Если он нашел пользователя и успешно просканировал, то он существует
        return 0, e.ErrTakenUsername
    }
    if !errors.Is(err, pgx.ErrNoRows) {
        // Если ошибка не является ошибкой "Не найден пользователь"
        return 0, err
    }

    tx, err := r.db.Begin(context.Background())

    // Создание пользователя
    err = tx.QueryRow(context.Background(), fmt.Sprintf(`
        INSERT INTO %s (username, name, email, phone, birthday, created_at)
        VALUES ($1, $2, $3, $4, $5, NOW())
        RETURNING id
    `, usersTable), req.Username, req.Name, req.Email, req.Phone, req.Birthday).Scan(&id)
    if err != nil {
        tx.Rollback(context.Background())
        return 0, err
    }

    // Привязка пароля
    err = r.db.QueryRow(context.Background(), fmt.Sprintf(`
        INSERT INTO %s (user_id, password)
        VALUES ($1, $2)
        RETURNING user_id
    `, authTable), id, req.Password).Scan(&id)
    if err != nil {
        tx.Rollback(context.Background())
        return 0, err
    }

    err = tx.Commit(context.Background())
    if err != nil {
        return 0, err
    }

    return id, nil
}

func (r *AuthRepository) isPasswordMatch(input *structs.AuthInput) (bool, error) {
    var userPassword string
    query := fmt.Sprintf(`SELECT user_id, password FROM %s WHERE user_id = $1`, authTable)
    err := r.db.QueryRow(context.Background(), query, *input.UserId).Scan(&userPassword)
    if err != nil {
        return false, err
    }

    if userPassword != *input.Password {
        return false, nil
    }

    return true, nil
}

func (r *AuthRepository) AddRefreshToken(input *structs.AuthInput, token string) error {
    isMatch, err := r.isPasswordMatch(input)
    if err != nil {
        return err
    }

    if !isMatch {
        return errors.New("невереный пароль")
    }

    var userId int
    query := fmt.Sprintf("INSERT INTO %s (user_id, token) values ($1, $2) RETURNING user_id", refreshTable)
    err = r.db.QueryRow(context.Background(), query, *input.UserId, token).Scan(&userId)
    return err
}
