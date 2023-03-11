package service

import (
	"context"
	"strconv"
	"time"

	"github.com/spliatz/bloggy-backend/internal/domain/entity"
	auth_helpers "github.com/spliatz/bloggy-backend/pkg/auth"
)

type authStorage interface {
	SetSession(ctx context.Context, userId int, session entity.Session) error
	DeleteUserSession(ctx context.Context, userId int) error
	CheckRefresh(ctx context.Context, refreshToken string) error
}

type authService struct {
	storage      authStorage
	tokenManager auth_helpers.TokenManager
}

func NewAuthService(storage authStorage, tokenManager auth_helpers.TokenManager) *authService {
	return &authService{storage: storage, tokenManager: tokenManager}
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
