package services

import (
    "context"
    "strconv"
    "time"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/auth"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/hash"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/utils"
    "github.com/jackc/pgx/v5/pgtype"
)

type AuthService struct {
    repos        *repository.Repository
    hasher       hash.PasswordHasher
    tokenManager auth.TokenManager
}

func newAuthService(
    repos *repository.Repository,
    hasher hash.PasswordHasher,
    tokenManager auth.TokenManager,
) *AuthService {
    return &AuthService{
        repos:        repos,
        hasher:       hasher,
        tokenManager: tokenManager,
    }
}

type SignUpInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`

    Name     *string `json:"name"`
    Birthday *string `json:"birthday"` // Use "2000-12-31" format
    Email    *string `json:"email"`
    Phone    *string `json:"phone"`
}

type Tokens struct {
    Access  string
    Refresh string
}

func (s *AuthService) SignUp(ctx context.Context, i SignUpInput) (Tokens, error) {

    if err := utils.CheckUsername(i.Username); err != nil {
        return Tokens{}, err
    }

    if err := utils.CheckPassword(i.Password); err != nil {
        return Tokens{}, err
    }

    passHash, err := s.hasher.Hash(i.Password)
    if err != nil {
        return Tokens{}, err
    }

    var name pgtype.Text
    if i.Name != nil && *i.Name != "" {
        name = pgtype.Text{String: *i.Name, Valid: true}
    }

    var birthday pgtype.Date
    if i.Birthday != nil && *i.Birthday != "" {
        date, err := utils.ParseDate(*i.Birthday)
        if err != nil {
            return Tokens{}, err
        }
        birthday = pgtype.Date{Time: date, Valid: true}
    }

    var email pgtype.Text
    if i.Email != nil && *i.Email != "" {
        email = pgtype.Text{String: *i.Email, Valid: true}
    }

    var phone pgtype.Text
    if i.Phone != nil && *i.Phone != "" {
        phone = pgtype.Text{String: *i.Phone, Valid: true}
    }

    userInput := repository.User{
        Username:  i.Username,
        Password:  passHash,
        Name:      name,
        Birthday:  birthday,
        Email:     email,
        Phone:     phone,
        CreatedAt: time.Now(),
    }

    user, err := s.repos.SignUp(ctx, userInput)
    if err != nil {
        return Tokens{}, err
    }

    return s.createSession(ctx, user.Id)
}

type SignInInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func (s *AuthService) SignIn(ctx context.Context, i SignInInput) (Tokens, error) {
    passHash, err := s.hasher.Hash(i.Password)
    if err != nil {
        return Tokens{}, err
    }

    user, err := s.repos.GetByCredentials(ctx, i.Username, passHash)
    if err != nil {
        return Tokens{}, err
    }

    return s.createSession(ctx, user.Id)
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
    /*
       1. Проверить токен
       2. Удалить в любом случае
       3. Если он просрочен, выдать ошибку
       4. Если не просрочен, создать новый токен
    */

    err := s.repos.CheckRefresh(ctx, refreshToken)
    if err != nil && !errors.Is(err, errors.ErrTokenExpired) {
        return Tokens{}, err
    }

    if err := s.repos.DeleteRefresh(ctx, refreshToken); err != nil {
        return Tokens{}, err
    }

    if errors.Is(err, errors.ErrTokenExpired) {
        return Tokens{}, errors.ErrTokenExpired
    }

    user, err := s.repos.GetByRefreshToken(ctx, refreshToken)
    if err != nil {
        return Tokens{}, nil
    }

    return s.createSession(ctx, user.Id)
}

func (s *AuthService) createSession(ctx context.Context, userId int) (Tokens, error) {
    var (
        res Tokens
        err error
    )

    res.Access, err = s.tokenManager.NewJWT(strconv.Itoa(userId), time.Minute*15)
    if err != nil {
        return Tokens{}, err
    }

    res.Refresh, err = s.tokenManager.NewRefreshToken()
    if err != nil {
        return Tokens{}, err
    }

    session := repository.Session{
        RefreshToken: res.Refresh,
        ExpiresAt:    time.Now().Add(time.Hour * 720), // 30 days
    }

    err = s.repos.SetSession(ctx, userId, session)
    if err != nil {
        return Tokens{}, err
    }

    return res, err
}
