package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/spliatz/bloggy-backend/internal/domain/entity"
	user_dto "github.com/spliatz/bloggy-backend/internal/domain/usecase/user/dto"
	"github.com/spliatz/bloggy-backend/pkg/errors"
	"github.com/spliatz/bloggy-backend/pkg/hash"
	"github.com/spliatz/bloggy-backend/pkg/utils"
)

type userStorage interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetByUsername(ctx context.Context, username string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error)
	EditById(ctx context.Context, id int, req map[string]string) (entity.UserResponse, error)
	EditNameById(ctx context.Context, id int, name string) (entity.UserResponse, error)
	EditBirthdayById(ctx context.Context, id int, birthday string) (entity.UserResponse, error)
	EditUsernameById(ctx context.Context, id int, username string) (entity.UserResponse, error)
	EditEmailById(ctx context.Context, id int, email string) (entity.UserResponse, error)
	EditPhoneById(ctx context.Context, id int, phone string) (entity.UserResponse, error)
}

type userCache interface {
	GetById(ctx context.Context, id int) (entity.User, error)
	GetByUsername(ctx context.Context, username string) (entity.User, error)
	Set(ctx context.Context, user entity.User)
}

type userService struct {
	storage userStorage
	hasher  hash.PasswordHasher
	cache   userCache
}

func NewUserService(storage userStorage, hasher hash.PasswordHasher, cache userCache) *userService {
	return &userService{storage: storage, hasher: hasher, cache: cache}
}

func (s *userService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	return s.cache.GetById(ctx, id)
	// return s.cache.GetById(ctx, id)
}

func (s *userService) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	return s.cache.GetByUsername(ctx, username)
}

func (s *userService) CreateUser(ctx context.Context, dto user_dto.CreateUserDTO) (int, error) {

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

	if id, err := s.storage.CreateUser(ctx, u); err != nil {
		return -1, err
	} else {
		u.Id = id
		s.cache.Set(ctx, u)
		return id, nil
	}
}

func (s *userService) GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error) {
	return s.storage.GetByRefreshToken(ctx, refreshToken)
}

func (s *userService) EditById(ctx context.Context, id int, i user_dto.EditUserDTO) (entity.UserResponse, error) {
	if id <= 0 {
		return entity.UserResponse{}, errors.ErrWrongId
	}

	reqMap := map[string]string{}
	if i.Username != nil {
		if err := utils.CheckUsername(*i.Username); err != nil {
			return entity.UserResponse{}, err
		}
		reqMap["username"] = *i.Username
	}
	if i.Name != nil && *i.Name != "" {
		if err := utils.CheckName(*i.Name); err != nil {
			return entity.UserResponse{}, err
		}
		reqMap["name"] = *i.Name
	}
	if i.Birthday != nil && *i.Birthday != "" {
		// TODO: Добавить проверку дня рождения на совпадение с форматом 2000-12-31
		reqMap["birthday"] = *i.Birthday
	}
	if i.Email != nil && *i.Email != "" {
		if err := utils.CheckEmail(*i.Email); err != nil {
			return entity.UserResponse{}, err
		}
		reqMap["email"] = *i.Email
	}
	if i.Phone != nil && *i.Phone != "" {
		// TODO: Сделать валидацию номера
		reqMap["phone"] = *i.Phone
	}

	user, err := s.storage.EditById(ctx, id, reqMap)
	if err != nil {
		return entity.UserResponse{}, err
	}

	return user, nil
}

func (s *userService) EditNameById(ctx context.Context, id int, name string) (entity.UserResponse, error) {
	return s.storage.EditNameById(ctx, id, name)
}

func (s *userService) EditBirthdayById(ctx context.Context, id int, birthday string) (entity.UserResponse, error) {
	return s.storage.EditBirthdayById(ctx, id, birthday)
}

func (s *userService) EditUsernameById(ctx context.Context, id int, username string) (entity.UserResponse, error) {
	return s.storage.EditUsernameById(ctx, id, username)
}

func (s *userService) EditEmailById(ctx context.Context, id int, email string) (entity.UserResponse, error) {
	if err := utils.CheckEmail(email); err != nil {
		return entity.UserResponse{}, err
	}
	return s.storage.EditEmailById(ctx, id, email)
}

func (s *userService) EditPhoneById(ctx context.Context, id int, phone string) (entity.UserResponse, error) {
	if err := utils.CheckPhone(phone); err != nil {
		return entity.UserResponse{}, err
	}
	return s.storage.EditPhoneById(ctx, id, phone)
}
