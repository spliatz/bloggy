package services

import (
	"github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
)

type AuthService struct {
	repos *repository.Repository
}

func newAuthService(repos *repository.Repository) *AuthService {
	return &AuthService{
		repos: repos,
	}
}

func (s *AuthService) Registration(input *structs.UserCreateInput) (int, error) {

	userId, err := s.repos.Registration(input)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (s *AuthService) Login() {

}

func (s *AuthService) generateAccessToken() {

}

func (s *AuthService) refreshAccessToken() {

}
