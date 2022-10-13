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
}

type Auth interface {
    SignUp(req *structs.SignUpRequest) (id int, err error)
}

func NewServices(repos *repository.Repository) *Services {
    return &Services{
        User: newUserService(repos),
        Auth: newAuthService(repos),
    }
}
