package user

import (
	"context"
	"github.com/spliatz/bloggy-backend/internal/domain/entity"
	user_dto "github.com/spliatz/bloggy-backend/internal/domain/usecase/user/dto"
)

type service interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetByUsername(ctx context.Context, username string) (entity.UserResponse, error)
	EditById(ctx context.Context, id int, i user_dto.EditUserDTO) (entity.UserResponse, error)
	EditNameById(ctx context.Context, id int, name string) (entity.UserResponse, error)
}

type postService interface {
	GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error)
}

type userUsecase struct {
	service
	postService
}

func NewUserUsecase(service service, postService postService) *userUsecase {
	return &userUsecase{service: service, postService: postService}
}

func (u *userUsecase) GetById(ctx context.Context, id int) (entity.UserResponse, error) {
	user, err := u.service.GetUserByID(ctx, id)
	return entity.UserResponse{
		Username:  user.Username,
		Name:      user.Name,
		Birthday:  user.Birthday,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}, err
}

func (u *userUsecase) GetByUsername(ctx context.Context, username string) (entity.UserResponse, error) {
	return u.service.GetByUsername(ctx, username)
}

func (u *userUsecase) GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error) {
	return u.postService.GetAllByUsername(ctx, username)
}

func (u *userUsecase) EditById(ctx context.Context, id int, dto user_dto.EditUserDTO) (entity.UserResponse, error) {
	return u.service.EditById(ctx, id, dto)
}

func (u *userUsecase) EditNameById(ctx context.Context, id int, dto user_dto.EditNameDTO) (entity.UserResponse, error) {
	return u.service.EditNameById(ctx, id, dto.Name)
}
