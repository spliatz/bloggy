package services

import (
	"github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
)

type Services struct {
	User
}

type User interface {
	Create(user structs.User) (*structs.User, error)
}

func NewServices(repos *repository.Repository) *Services {
	return &Services{
		User: newUserService(repos),
	}
}
