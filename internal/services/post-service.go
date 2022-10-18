package services

import (
    "time"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
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

func (s *PostService) Create(req CreatePostInput) (int, error) {
    post := repository.Post{
        Content:   req.Content,
        CreatedAt: time.Now(),
        UserId:    req.AuthorId,
    }

    id, err := s.repos.PostRepo.Create(post)
    if err != nil {
        return 0, err
    }

    return id, nil
}

func (s *PostService) GetById(id int) (repository.Post, error) {
    return s.repos.PostRepo.GetById(id)
}

func (s *PostService) GetAllByUsername(username string) ([]repository.Post, error) {
    return s.repos.GetAllByUsername(username)
}

func (s *PostService) DeleteById(postId int) error {
    return s.repos.PostRepo.DeleteById(postId)
}
