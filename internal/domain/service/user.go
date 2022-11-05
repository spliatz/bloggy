package service

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgtype"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    user_usecase "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/user"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/hash"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/utils"
)

type UserStorage interface {
    GetUserByID(ctx context.Context, id int) (entity.User, error)
    CreateUser(ctx context.Context, user entity.User) (int, error)
    GetUserByUserNameAndPassword(ctx context.Context, username, password string) (entity.User, error)
    GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error)
}

type userService struct {
    storage UserStorage
    hasher  hash.PasswordHasher
}

func NewUserService(storage UserStorage, hasher hash.PasswordHasher) *userService {
    return &userService{storage: storage, hasher: hasher}
}

func (s *userService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
    return s.storage.GetUserByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, dto user_usecase.CreateUserDTO) (int, error) {

    if err := utils.CheckUsername(dto.Username); err != nil {
        return -1, err
    }

    if err := utils.CheckPassword(dto.Password); err != nil {
        return -1, err
    }

    var name pgtype.Text
    if dto.Name != nil && *dto.Name != "" {
        name = pgtype.Text{String: *dto.Name, Valid: true}
    }

    var birthday pgtype.Date
    if dto.Birthday != nil && *dto.Birthday != "" {
        date, err := utils.ParseDate(*dto.Birthday)
        if err != nil {
            return -1, err
        }
        birthday = pgtype.Date{Time: date, Valid: true}
    }

    var email pgtype.Text
    if dto.Email != nil && *dto.Email != "" {
        email = pgtype.Text{String: *dto.Email, Valid: true}
    }

    var phone pgtype.Text
    if dto.Phone != nil && *dto.Phone != "" {
        phone = pgtype.Text{String: *dto.Phone, Valid: true}
    }

    passHash, err := s.hasher.Hash(dto.Password)
    if err != nil {
        return -1, err

    }

    u := entity.User{
        Username:  dto.Username,
        Password:  passHash,
        Name:      name,
        Birthday:  birthday,
        Email:     email,
        Phone:     phone,
        CreatedAt: time.Now(),
    }

    return s.storage.CreateUser(ctx, u)
}

func (s *userService) GetByCredentials(ctx context.Context, dto user_usecase.GetByCredentialsDTO) (entity.User, error) {
    passHash, err := s.hasher.Hash(dto.Password)
    if err != nil {
        return entity.User{}, err
    }

    return s.storage.GetUserByUserNameAndPassword(ctx, dto.Username, passHash)
}

func (s *userService) GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error) {
    return s.storage.GetByRefreshToken(ctx, refreshToken)
}
