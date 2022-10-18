package repository

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5"

    e "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
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
    Id        int       `db:"id"`
    AuthorId  int       `db:"author_id"`
    Content   string    `db:"content"`
    CreatedAt time.Time `db:"created_at"`
}

func (r *PostRepository) Create(c context.Context, p Post) (id int, err error) {

    err = r.db.QueryRow(c, fmt.Sprintf(`
        INSERT INTO %s(author_id, content, created_at)
        VALUES ($1, $2, $3)
        RETURNING id
    `, postsTable), p.AuthorId, p.Content, p.CreatedAt).Scan(&id)
    if err != nil {
        return 0, err
    }

    return id, nil
}

func (r *PostRepository) GetById(c context.Context, id int) (p Post, err error) {
    err = r.db.QueryRow(c, fmt.Sprintf(`
        SELECT id, author_id, content, created_at
        FROM %s
        WHERE id = $1
    `, postsTable), id).Scan(&p.Id, &p.AuthorId, &p.Content, &p.CreatedAt)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return Post{}, e.ErrIdNotFound
        }

        return Post{}, err
    }

    return p, nil
}

func (r *PostRepository) GetAllByUsername(c context.Context, username string) (posts []Post, err error) {
    var authorId int
    err = r.db.QueryRow(c, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE username = $1
    `, usersTable), username).Scan(&authorId)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return make([]Post, 0), e.ErrUsernameNotFound
        }

        return make([]Post, 0), err
    }

    rows, err := r.db.Query(c, fmt.Sprintf(`
        SELECT id, author_id, content, created_at 
        FROM %s WHERE author_id = $1
    `, postsTable), authorId)
    if err != nil {
        return make([]Post, 0), err
    }

    posts = make([]Post, 0)
    for rows.Next() {
        var p Post
        err = rows.Scan(&p.Id, &p.AuthorId, &p.Content, &p.CreatedAt)
        if err != nil {
            return make([]Post, 0), err
        }
        posts = append(posts, p)
    }

    return posts, nil
}

func (r *PostRepository) DeleteById(c context.Context, id int) error {
    var postId int
    err := r.db.QueryRow(c, fmt.Sprintf(`
        DELETE FROM %s
        WHERE id = $1
        RETURNING id
    `, postsTable), id).Scan(&postId)
    if errors.Is(err, pgx.ErrNoRows) {
        return e.ErrIdNotFound
    }

    return err
}

func (r *PostRepository) IsAuthor(c context.Context, postId int, authorId int) (bool, error) {
    p := Post{}
    err := r.db.QueryRow(c, fmt.Sprintf(`
        SELECT id, author_id, content, created_at
        FROM %s
        WHERE id = $1 AND author_id = $2
    `, postsTable), postId, authorId).Scan(&p.Id, &p.AuthorId, &p.Content, &p.CreatedAt)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return false, nil
        }

        return false, err
    }

    return true, nil
}
