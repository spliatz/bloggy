package services

import (
    "context"
    "strings"
    "time"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/utils"
)

type UserService struct {
    repos *repository.Repository
}

func newUserService(repos *repository.Repository) *UserService {
    return &UserService{
        repos: repos,
    }
}

type UserResponse struct {
    Username  string    `json:"username"`
    Name      *string   `json:"name"`
    Birthday  *string   `json:"birthday"`
    Email     *string   `json:"email"`
    Phone     *string   `json:"phone"`
    CreatedAt time.Time `json:"created_at"`
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (UserResponse, error) {
    if err := utils.CheckUsername(username); err != nil {
        return UserResponse{}, err
    }

    user, err := s.repos.GetByUsername(ctx, username)
    if err != nil {
        return UserResponse{}, err
    }

    return s.uToUres(user), nil
}

type EditInput struct {
    Username *string `json:"username"`
    Name     *string `json:"name"`
    Birthday *string `json:"birthday"`
    Email    *string `json:"email"`
    Phone    *string `json:"phone"`
}

func (s *UserService) EditById(ctx context.Context, id int, i EditInput) (UserResponse, error) {
    if id <= 0 {
        return UserResponse{}, errors.ErrWrongId
    }

    reqMap := map[string]string{}
    if i.Username != nil {
        if err := utils.CheckUsername(*i.Username); err != nil {
            return UserResponse{}, err
        }
        reqMap["username"] = *i.Username
    }
    if i.Name != nil && *i.Name != "" {
        if err := utils.CheckName(*i.Name); err != nil {
            return UserResponse{}, err
        }
        reqMap["name"] = *i.Name
    }
    if i.Birthday != nil && *i.Birthday != "" {
        // TODO: Добавить проверку дня рождения на совпадение с форматом 2000-12-31
        reqMap["birthday"] = *i.Birthday
    }
    if i.Email != nil && *i.Email != "" {
        // TODO: Добавить проверку почты на совпадение с форматом xxx@zzz.com
        reqMap["email"] = *i.Email
    }
    if i.Phone != nil && *i.Phone != "" {
        // TODO: Сделать валидацию номера
        reqMap["phone"] = *i.Phone
    }

    user, err := s.repos.EditById(ctx, id, reqMap)
    if err != nil {
        return UserResponse{}, err
    }

    return s.uToUres(user), nil
}

func (s *UserService) uToUres(user repository.User) UserResponse {
    newUser := UserResponse{
        Username:  user.Username,
        CreatedAt: user.CreatedAt,
    }

    if user.Name.Valid {
        newUser.Name = &user.Name.String
    }
    if user.Birthday.Valid {
        dateWithTime := user.Birthday.Time.Format(time.RFC3339)
        dateOnly, _, _ := strings.Cut(dateWithTime, "T")
        newUser.Birthday = &dateOnly
    }
    if user.Email.Valid {
        newUser.Email = &user.Email.String
    }
    if user.Phone.Valid {
        newUser.Phone = &user.Phone.String
    }

    return newUser
}
