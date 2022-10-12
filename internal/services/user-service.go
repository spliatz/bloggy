package services

import (
    "errors"
    "github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
    "github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
)

type UserService struct {
    repos *repository.Repository
}

func newUserService(repos *repository.Repository) *UserService {
    return &UserService{
        repos: repos,
    }
}

func (s *UserService) Create(input *structs.UserCreateInput) (*structs.User, error) {

    if input.Username == nil {
        return nil, errors.New("поле username обязательно")
    }

    user, err := s.repos.User.Create(input)
    if err != nil {
        return nil, err
    }

    return user, nil
}
