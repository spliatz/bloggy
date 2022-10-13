package structs

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type User struct {
	Id        int         `json:"id" db:"id"`
	Username  string      `json:"username" db:"username" binding:"required"`
	Name      pgtype.Text `json:"name" db:"name"`
	Birthday  pgtype.Date `json:"birthday" db:"birthday"`
	Email     pgtype.Text `json:"email" db:"email"`
	Phone     pgtype.Text `json:"phone" db:"phone"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
}

type UserCreateInput struct {
	Username *string `json:"username" db:"username"`
	Name     *string `json:"name" db:"name" binding:"required"`
	Birthday *string `json:"birthday" db:"birthday"`
	Email    *string `json:"email" db:"email"`
	Phone    *string `json:"phone" db:"phone"`
	Password *string `json:"password" db:"password" binding:"required"`
}
