package service

import (
    "context"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
)

type postStorage interface {
    GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error)
}

type postService struct {
    postStorage
}

func NewPostService(postStorage postStorage) *postService {
    return &postService{postStorage: postStorage}
}

func (s *postService) GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error) {
    return s.postStorage.GetAllByUsername(ctx, username)
}
