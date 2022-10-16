package repository

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5"

    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type PostRepository struct {
    db *pgx.Conn
}

func newPostRepository(db *pgx.Conn) *PostRepository {
    return &PostRepository{
        db: db,
    }
}

type Post struct {
    Id        int       `json:"id" db:"id"`
    UserId    int       `json:"author_id" db:"author_id"`
    Content   string    `json:"content" db:"content"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (r *PostRepository) Create(req Post) (int, error) {
    var postId int

    err := r.db.QueryRow(context.Background(), fmt.Sprintf(
        `INSERT INTO %s (author_id, content, created_at) values ($1, $2, $3) RETURNING id`, postsTable),
        req.UserId, req.Content, req.CreatedAt).Scan(&postId)
    if err != nil {
        return 0, err
    }

    return postId, nil
}

func (r *PostRepository) GetOne(id int) (Post, error) {
    var post Post
    err := r.db.QueryRow(context.Background(),
        fmt.Sprintf(`SELECT id, author_id, content, created_at FROM %s WHERE id = $1`, postsTable),
        id).Scan(&post.Id, &post.UserId, &post.Content, &post.CreatedAt)

    return post, err
}

func (r *PostRepository) GetAllUserPosts(username string) ([]Post, error) {

    // получаем юзера
    var userId int
    err := r.db.QueryRow(context.Background(), fmt.Sprintf(
        `SELECT id FROM %s WHERE username = $1`, usersTable), username).Scan(&userId)
    if err != nil {
        return nil, errors.ErrUserDoesNotExist
    }

    // получаем все посты полученного юзера
    rows, err := r.db.Query(context.Background(), fmt.Sprintf(
        `
                SELECT id, author_id, content, created_at 
                FROM %s WHERE author_id = u.id`,
        postsTable), userId)

    if err != nil {
        return nil, err
    }

    posts := make([]Post, 0)

    // запись всех пришедших с базы данных постов в массив posts
    for rows.Next() {
        var post Post
        err = rows.Scan(&post.Id, &post.UserId, &post.Content, &post.CreatedAt)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}

func (r *PostRepository) Delete(id int) error {
    var postId int
    err := r.db.QueryRow(context.Background(),
        fmt.Sprintf(`DELETE FROM %s WHERE id = $1 RETURNING id`, postsTable),
        id).Scan(&postId)
    return err
}
