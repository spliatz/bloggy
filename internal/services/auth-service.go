package services

import (
    "github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
    "github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/utils"
)

type AuthService struct {
    repos *repository.Repository
}

func newAuthService(repos *repository.Repository) *AuthService {
    return &AuthService{
        repos: repos,
    }
}

func (s *AuthService) SignUp(req *structs.SignUpRequest) (id int, err error) {

    if err = utils.CheckUsername(req.Username); err != nil {
        return 0, err
    }

    if err = utils.CheckPassword(req.Password); err != nil {
        return 0, err
    }

    id, err = s.repos.SignUp(req)
    if err != nil {
        return 0, err
    }

    return id, nil
}

func (s *AuthService) generateAccessToken() {

}

func (s *AuthService) refreshAccessToken() {

}
