package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"
)

type AuthRepository struct {
	db *pgx.Conn
}

func newAuthRepository(db *pgx.Conn) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Registration(input *structs.UserCreateInput) (int, error) {

	args, argsValue, argsOrder, err := r.prepareCreateUserQuery(input)
	if err != nil {
		return -1, err
	}

	tx, err := r.db.Begin(context.Background())

	createUserQueryString := fmt.Sprintf(`
        INSERT INTO %s (%s)
        values(%s)
        RETURNING id, name, username, created_at, birthday, email, phone
    `, usersTable, strings.Join(args, ", "), strings.Join(argsOrder, ", "))

	user := structs.User{}
	err = tx.QueryRow(context.Background(), createUserQueryString, argsValue...).Scan(
		&user.Id, &user.Name, &user.Username, &user.CreatedAt, &user.Birthday, &user.Email, &user.Phone)

	if err != nil {
		tx.Rollback(context.Background())
		return -1, err
	}

	var userId int
	query := fmt.Sprintf(`INSERT INTO %s (user_id, password) values ($1, $2) RETURNING user_id`, authTable)
	err = r.db.QueryRow(context.Background(), query, user.Id, *input.Password).Scan(&userId)
	if err != nil {
		tx.Rollback(context.Background())
		return -1, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return -1, err
	}

	return userId, nil
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

func (r *AuthRepository) prepareCreateUserQuery(input *structs.UserCreateInput) (
	args []string, argsValue []interface{}, argsOrder []string, err error) {

	args = make([]string, 0)
	argsValue = make([]interface{}, 0)
	argsOrder = make([]string, 0)
	argId := 1

	args = append(args, "username")
	argsValue = append(argsValue, *input.Username)
	argsOrder = append(argsOrder, fmt.Sprintf("$%x", argId))
	argId++

	args = append(args, "name")
	argsValue = append(argsValue, *input.Name)
	argsOrder = append(argsOrder, fmt.Sprintf("$%x", argId))
	argId++
	args = append(args, "created_at")
	argsValue = append(argsValue, time.Now())
	argsOrder = append(argsOrder, fmt.Sprintf("$%x", argId))
	argId++

	if input.Birthday != nil {
		birthday, err := time.Parse("02-01-2006", *input.Birthday)
		if err != nil {
			return nil, nil, nil, err
		}

		args = append(args, "birthday")
		argsValue = append(argsValue, birthday)
		argsOrder = append(argsOrder, fmt.Sprintf("$%x", argId))
		argId++
	}

	if input.Email != nil {
		args = append(args, fmt.Sprintf("email"))
		argsValue = append(argsValue, *input.Email)
		argsOrder = append(argsOrder, fmt.Sprintf("$%x", argId))
		argId++
	}

	if input.Phone != nil {
		args = append(args, fmt.Sprintf("phone"))
		argsValue = append(argsValue, *input.Phone)
		argsOrder = append(argsOrder, fmt.Sprintf("$%x", argId))
		argId++
	}

	return
}
