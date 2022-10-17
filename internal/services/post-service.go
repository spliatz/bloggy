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
    Content string `json:"content" binding:"required"`
    UserId  int    `json:"author_id"`
}

func (s *PostService) Create(req CreatePostInput) (int, error) {
    createdAt := time.Now()
    post := repository.Post{
        Content:   req.Content,
        CreatedAt: createdAt,
        UserId:    req.UserId,
    }
    id, err := s.repos.PostRepo.Create(post)
    if err != nil {
        return 0, err
    }

    return id, nil
}

func (s *PostService) GetOneById(id int) (repository.Post, error) {
    return s.repos.PostRepo.GetOneById(id)
}

func (s *PostService) GetAllUserPosts(username string) ([]repository.Post, error) {
    return s.repos.GetAllUserPosts(username)
}

func (s *PostService) DeleteById(postId int) error {
    return s.repos.PostRepo.DeleteById(postId)
}
