package services

import (
	"github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
)

type Services struct {
	User
	Auth
}

type User interface {
	Create(input *structs.UserCreateInput) (*structs.User, error)
}

type Auth interface {
	Registration(input *structs.UserCreateInput) (int, error)
}

func NewServices(repos *repository.Repository) *Services {
	return &Services{
		User: newUserService(repos),
		Auth: newAuthService(repos),
	}
}
