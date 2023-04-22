package service

import (
	"context"
	e "github.com/spliatz/bloggy-backend/pkg/errors"
	"github.com/spliatz/bloggy-backend/pkg/hash"
	"strconv"
	"time"

	"github.com/spliatz/bloggy-backend/internal/domain/entity"
	auth_helpers "github.com/spliatz/bloggy-backend/pkg/auth"
)

type authStorage interface {
	SetSession(ctx context.Context, userId int, session entity.Session) error
	DeleteUserSession(ctx context.Context, userId int) error
	CheckRefresh(ctx context.Context, refreshToken string) error
	UpdatePassword(ctx context.Context, userId int, newPassword string) error
	GetPassword(ctx context.Context, userId int) (string, error)
}

type authService struct {
	storage      authStorage
	tokenManager auth_helpers.TokenManager
	hasher       hash.PasswordHasher
}

func NewAuthService(storage authStorage, tokenManager auth_helpers.TokenManager, hasher hash.PasswordHasher) *authService {
	return &authService{storage: storage, tokenManager: tokenManager, hasher: hasher}
}

func (s *authService) GenerateAccessToken(ctx context.Context, userId int) (string, error) {
	return s.tokenManager.NewJWT(strconv.Itoa(userId), time.Minute*15)
}

func (s *authService) GenerateRefreshToken(ctx context.Context) (string, error) {
	return s.tokenManager.NewRefreshToken()
}

func (s *authService) SetSession(ctx context.Context, userId int, session entity.Session) error {
	return s.storage.SetSession(ctx, userId, session)
}

func (s *authService) CheckRefresh(ctx context.Context, refreshToken string) error {
	return s.storage.CheckRefresh(ctx, refreshToken)
}

func (s *authService) DeleteUserSession(ctx context.Context, userId int) error {
	return s.storage.DeleteUserSession(ctx, userId)
}

func (s *authService) CheckPassword(ctx context.Context, userId int, password string) error {
	storPass, err := s.storage.GetPassword(ctx, userId)

	if err != nil {
		return err
	}

	hashed, err := s.EncryptPassword(ctx, password)
	if err != nil {
		return err

	}

	if storPass != hashed {
		return e.ErrWrongCredentialsPassword
	}

	return nil
}

func (s *authService) UpdatePassword(ctx context.Context, userId int, newPassword string) error {
	hashed, err := s.EncryptPassword(ctx, newPassword)
	if err != nil {
		return err
	}
	return s.storage.UpdatePassword(ctx, userId, hashed)
}

func (s *authService) EncryptPassword(ctx context.Context, password string) (string, error) {
	hashed, err := s.hasher.Hash(password)
	if err != nil {
		return "", err
	}
	return hashed, nil
}
