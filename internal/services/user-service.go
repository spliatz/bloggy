package services

import (
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

func (s *UserService) Create(u structs.User) (*structs.User, error) {
    err := s.repos.User.Create(u.Username, u.Name, u.Surname, u.Email, u.Password)
    if err != nil {
        return nil, err
    }

    user := &structs.User{
        Username: u.Username,
        Name:     u.Name,
        Surname:  u.Surname,
        Email:    u.Email,
        Password: u.Password,
    }

    return user, nil
}
