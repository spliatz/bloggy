package services

import (
    "context"
    "time"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/utils"
)

type PostService struct {
    repos *repository.Repository
}

func newPostService(repos *repository.Repository) *PostService {
    return &PostService{
        repos: repos,
    }
}

type CreatePostInput struct {
    Content  string `json:"content" binding:"required"`
    AuthorId int    `json:"author_id"`
}

func (s *PostService) Create(c context.Context, i CreatePostInput) (int, error) {
    post := repository.Post{
        Content:   i.Content,
        CreatedAt: time.Now(),
        AuthorId:  i.AuthorId,
    }

    id, err := s.repos.PostRepo.Create(c, post)
    if err != nil {
        return 0, err
    }

    return id, nil
}

type PostResponse struct {
    Id             int     `json:"id"`
    AuthorUsername string  `json:"author_username"`
    AuthorName     *string `json:"author_name"`
    Content        string  `json:"content"`
    CreatedAt      int     `json:"created_at"`
}

func (s *PostService) GetById(c context.Context, id int) (PostResponse, error) {
    post, err := s.repos.PostRepo.GetById(c, id)
    if err != nil {
        return PostResponse{}, err
    }

    author, err := s.repos.UserRepo.GetById(c, post.AuthorId)
    if err != nil {
        return PostResponse{}, err
    }

    pr := s.PtoPr(post, author)

    return pr, nil
}

func (s *PostService) GetAllByUsername(c context.Context, username string) ([]PostResponse, error) {
    err := utils.CheckUsername(username)
    if err != nil {
        return make([]PostResponse, 0), err
    }

    posts, err := s.repos.PostRepo.GetAllByUsername(c, username)
    if err != nil {
        return make([]PostResponse, 0), err
    }

    author, err := s.repos.UserRepo.GetByUsername(c, username)
    if err != nil {
        return make([]PostResponse, 0), err
    }

    resPosts := make([]PostResponse, 0)
    for _, p := range posts {
        resPosts = append(resPosts, s.PtoPr(p, author))
    }

    return resPosts, nil
}

func (s *PostService) DeleteById(c context.Context, postId int) error {
    return s.repos.PostRepo.DeleteById(c, postId)
}

func (s *PostService) IsAuthor(c context.Context, postId int, authorId int) (bool, error) {
    return s.repos.PostRepo.IsAuthor(c, postId, authorId)
}

func (s *PostService) PtoPr(p repository.Post, a repository.User) PostResponse {
    pr := PostResponse{
        Id:             p.Id,
        AuthorUsername: a.Username,
        Content:        p.Content,
        CreatedAt:      int(p.CreatedAt.Unix()),
    }

    if a.Name.Valid {
        pr.AuthorName = &a.Name.String
    }

    return pr
}
