package repository

import (
    "context"
    "fmt"
    "strings"
    "time"

    "github.com/jackc/pgx/v5"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
)

type UserRepository struct {
    db *pgx.Conn
}

func newUserRepository(db *pgx.Conn) *UserRepository {
    return &UserRepository{
        db: db,
    }
}

func (r *UserRepository) Create(input *structs.UserCreateInput) (*structs.User, error) {

    args := make([]string, 0)
    argsValue := make([]interface{}, 0)
    argsOrder := make([]string, 0)
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
            return nil, err
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

    queryString := fmt.Sprintf(`
        INSERT INTO %s (%s)
        values(%s)
        RETURNING id, name, username, created_at, birthday, email, phone
    `, usersTable, strings.Join(args, ", "), strings.Join(argsOrder, ", "))

    user := structs.User{}
    err := r.db.QueryRow(context.Background(), queryString, argsValue...).Scan(
        &user.Id,
        &user.Name,
        &user.Username,
        &user.CreatedAt,
        &user.Birthday,
        &user.Email,
        &user.Phone)

    if err != nil {
        return nil, err
    }

    return &user, nil
}
