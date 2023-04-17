package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/spliatz/bloggy-backend/internal/domain/entity"
	e "github.com/spliatz/bloggy-backend/pkg/errors"
)

type authStorage struct {
	db *pgx.Conn
}

func NewAuthStorage(db *pgx.Conn) *authStorage {
	return &authStorage{db: db}
}

func (s *authStorage) SetSession(ctx context.Context, userId int, session entity.Session) error {
	_, err := s.db.Exec(ctx, fmt.Sprintf(`
        INSERT INTO %s(user_id, token, expires_at)
        VALUES ($1, $2, $3)
    `, refreshTable), userId, session.RefreshToken, session.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *authStorage) DeleteUserSession(ctx context.Context, userId int) error {
	err := s.db.QueryRow(ctx, fmt.Sprintf(`
                DELETE FROM %s
                WHERE user_id = $1
                RETURNING user_id
            `, refreshTable), userId).Scan(&userId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return e.ErrTokenNotFound
	}

	return nil
}

func (s *authStorage) CheckRefresh(ctx context.Context, refreshToken string) error {
	var expiresAt pgtype.Date
	err := s.db.QueryRow(ctx, fmt.Sprintf(`
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

func (s *authStorage) GetPassword(ctx context.Context, userId int) (string, error) {
	var password string

	query := fmt.Sprintf(`SELECT password FROM %s WHERE user_id=$1`, authTable)
	if err := s.db.QueryRow(ctx, query, userId).Scan(&password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", e.ErrWrongCredentialsPassword
		}

		return "", err
	}

	return password, nil
}

func (s *authStorage) UpdatePassword(ctx context.Context, userId int, newPassword string) error {
	var password string
	query := fmt.Sprintf(`UPDATE %s SET password=$1 WHERE user_id=$2 RETURNING password`, authTable)
	if err := s.db.QueryRow(ctx, query, newPassword, userId).Scan(&password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.ErrUserNotFound
		}

		return err
	}

	return nil
}
