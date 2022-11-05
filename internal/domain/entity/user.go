package entity

import (
    "time"

    "github.com/jackc/pgx/v5/pgtype"
)

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
