package postgres

import (
    "context"
    "errors"
    "fmt"

    "github.com/jackc/pgx/v5"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    e "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type postStorage struct {
    db *pgx.Conn
}

func NewPostStorage(db *pgx.Conn) *postStorage {
    return &postStorage{db: db}
}

func (s *postStorage) GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error) {
    var authorId int
    err = s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE username = $1
    `, usersTable), username).Scan(&authorId)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return make([]entity.Post, 0), e.ErrUsernameNotFound
        }

        return make([]entity.Post, 0), err
    }

    rows, err := s.db.Query(ctx, fmt.Sprintf(`
        SELECT id, author_id, content, created_at 
        FROM %s WHERE author_id = $1
    `, postsTable), authorId)
    if err != nil {
        return make([]entity.Post, 0), err
    }

    posts = make([]entity.Post, 0)
    for rows.Next() {
        var p entity.Post
        err = rows.Scan(&p.Id, &p.AuthorId, &p.Content, &p.CreatedAt)
        if err != nil {
            return make([]entity.Post, 0), err
        }
        posts = append(posts, p)
    }

    return posts, nil
}
