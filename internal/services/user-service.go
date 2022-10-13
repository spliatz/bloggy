package services

import (
    "github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
)

type UserService struct {
    repos *repository.Repository
}

func newUserService(repos *repository.Repository) *UserService {
    return &UserService{
        repos: repos,
    }
}
