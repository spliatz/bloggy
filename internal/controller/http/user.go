package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/spliatz/bloggy-backend/internal/controller/http/response"
	"github.com/spliatz/bloggy-backend/internal/domain/entity"
	user_dto "github.com/spliatz/bloggy-backend/internal/domain/usecase/user/dto"
	"github.com/spliatz/bloggy-backend/pkg/errors"
)

type authMiddleware interface {
	UserIdentity(c *gin.Context)
}

type userUsecase interface {
	GetById(ctx context.Context, id int) (entity.UserResponse, error)
	GetByUsername(ctx context.Context, dto user_dto.GetByUsernameDTO) (entity.UserResponse, error)
	EditById(ctx context.Context, id int, dto user_dto.EditUserDTO) (entity.UserResponse, error)
	EditNameById(ctx context.Context, id int, dto user_dto.EditNameDTO) (entity.UserResponse, error)
	EditBirthdayById(ctx context.Context, id int, dto user_dto.EditBirthdayDTO) (entity.UserResponse, error)
	EditUsernameById(ctx context.Context, id int, dto user_dto.EditUsernameDTO) (entity.UserResponse, error)
	EditEmailById(ctx context.Context, id int, dto user_dto.EditEmailDTO) (entity.UserResponse, error)
	EditPhoneById(ctx context.Context, id int, dto user_dto.EditPhoneDTO) (entity.UserResponse, error)
	GetAllByUsername(ctx context.Context, dto user_dto.GetAllByUsernameDTO) (posts []entity.Post, err error)
}

type userHandler struct {
	authMiddleware
	userUsecase
}

func NewUserHandler(authMiddleware authMiddleware, userUsecase userUsecase) *userHandler {
	return &userHandler{authMiddleware: authMiddleware, userUsecase: userUsecase}
}

func (h *userHandler) Register(router *gin.Engine) {
	user := router.Group("/user")
	{
		protected := user.Group("", h.authMiddleware.UserIdentity)
		{
			protected.GET("/my", h.getMy)
			protected.PATCH("", h.editById)
			protected.PATCH("/name", h.editNameById)
			protected.PATCH("/birthday", h.editBirthdayById)
			protected.PATCH("/username", h.editUsernameById)
			protected.PATCH("/email", h.editEmailById)
			protected.PATCH("/phone", h.editPhoneById)
		}
		user.GET("/:username", h.getByUsername)
		user.GET("/:username/posts", h.getAllByUsername)
	}
}

// @Summary GetUserByUsername
// @Tags user
// @Description get my user info
// @Security ApiKeyAuth
// @ID get-my
// @Accept json
// @Produce json
// @Success 200 {object} entity.UserResponseSwagger
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /user/my [get]
func (h *userHandler) getMy(c *gin.Context) {
	userIdI, exist := c.Get(fieldUserId)
	if !exist {
		response.ResponseWithError(c, errors.ErrIdNotFound)
		return
	}

	userId, _ := userIdI.(int)
	user, err := h.userUsecase.GetById(c.Request.Context(), userId)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary GetUserByUsername
// @Tags user
// @Description get user by username
// @ID get-user-by-username
// @Accept json
// @Produce json
// @Param username path string true "User username"
// @Success 200 {object} entity.UserResponseSwagger
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /user/{username} [get]
func (h *userHandler) getByUsername(c *gin.Context) {
	username := c.Param(paramUsername)
	dto := user_dto.GetByUsernameDTO{Username: username}
	user, err := h.userUsecase.GetByUsername(c.Request.Context(), dto)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary EditUserByUsername
// @Tags user
// @Description login
// @Security ApiKeyAuth
// @ID edit-user-by-username
// @Accept json
// @Produce json
// @Param input body dto.EditUserDTO true "user information"
// @Success 200 {object} entity.UserResponseSwagger
// @Failure 400,401,403,404,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /user [patch]
func (h *userHandler) editById(c *gin.Context) {
	userIdI, exist := c.Get(fieldUserId)
	if !exist {
		response.ResponseWithError(c, errors.ErrIdNotFound)
		return
	}

	userId, _ := userIdI.(int)

	i := user_dto.EditUserDTO{}
	if err := c.BindJSON(&i); err != nil {
		response.ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
		return
	}

	user, err := h.userUsecase.EditById(c.Request.Context(), userId, i)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary EditName
// @Tags user
// @Description Edit user's name
// @Security ApiKeyAuth
// @ID edit-user-name
// @Accept json
// @Produce json
// @Param input body dto.EditNameDTO true "user name"
// @Success 200 {object} entity.UserResponseSwagger
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /user/name [patch]
func (h *userHandler) editNameById(c *gin.Context) {
	userIdI, exist := c.Get(fieldUserId)
	if !exist {
		response.ResponseWithError(c, errors.ErrIdNotFound)
		return
	}

	userId, _ := userIdI.(int)

	dto := user_dto.EditNameDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
		return
	}

	user, err := h.EditNameById(c.Request.Context(), userId, dto)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary EditUsername
// @Tags user
// @Description Edit user's username
// @Security ApiKeyAuth
// @ID edit-user-username
// @Accept json
// @Produce json
// @Param input body dto.EditUsernameDTO true "user username"
// @Success 200 {object} entity.UserResponseSwagger
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /user/username [patch]
func (h *userHandler) editUsernameById(c *gin.Context) {
	userIdI, exist := c.Get(fieldUserId)
	if !exist {
		response.ResponseWithError(c, errors.ErrIdNotFound)
		return
	}

	userId, _ := userIdI.(int)

	dto := user_dto.EditUsernameDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
		return
	}

	user, err := h.EditUsernameById(c.Request.Context(), userId, dto)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary EditEMail
// @Tags user
// @Description Edit user's email
// @Security ApiKeyAuth
// @ID edit-user-email
// @Accept json
// @Produce json
// @Param input body dto.EditEmailDTO true "user email"
// @Success 200 {object} entity.UserResponseSwagger
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /user/email [patch]
func (h *userHandler) editEmailById(c *gin.Context) {
	userIdI, exist := c.Get(fieldUserId)
	if !exist {
		response.ResponseWithError(c, errors.ErrIdNotFound)
		return
	}

	userId, _ := userIdI.(int)

	dto := user_dto.EditEmailDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
		return
	}

	user, err := h.EditEmailById(c.Request.Context(), userId, dto)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary EditPhone
// @Tags user
// @Description Edit user's phone
// @Security ApiKeyAuth
// @ID edit-user-phone
// @Accept json
// @Produce json
// @Param input body dto.EditPhoneDTO true "user phone"
// @Success 200 {object} entity.UserResponseSwagger
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /user/phone [patch]
func (h *userHandler) editPhoneById(c *gin.Context) {
	userIdI, exist := c.Get(fieldUserId)
	if !exist {
		response.ResponseWithError(c, errors.ErrIdNotFound)
		return
	}

	userId, _ := userIdI.(int)

	dto := user_dto.EditPhoneDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
		return
	}

	user, err := h.EditPhoneById(c.Request.Context(), userId, dto)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary EditBirthday
// @Tags user
// @Description Edit user's birthday
// @Security ApiKeyAuth
// @ID edit-user-birthday
// @Accept json
// @Produce json
// @Param input body dto.EditBirthdayDTO true "user birthday"
// @Success 200 {object} entity.UserResponseSwagger
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /user/birthday [patch]
func (h *userHandler) editBirthdayById(c *gin.Context) {
	userIdI, exist := c.Get(fieldUserId)
	if !exist {
		response.ResponseWithError(c, errors.ErrIdNotFound)
		return
	}

	userId, _ := userIdI.(int)

	dto := user_dto.EditBirthdayDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
		return
	}

	user, err := h.userUsecase.EditBirthdayById(c.Request.Context(), userId, dto)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary GetAllByUsername
// @Tags user
// @Description Get All User's Posts
// @ID get-all-by-username
// @Accept json
// @Produce json
// @Param username path string true "User username"
// @Success 200 {array} entity.Post
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /user/{username}/posts [get]
func (h *userHandler) getAllByUsername(c *gin.Context) {
	username := c.Param(paramUsername)
	dto := user_dto.GetAllByUsernameDTO{Username: username}
	posts, err := h.userUsecase.GetAllByUsername(c.Request.Context(), dto)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.JSON(http.StatusOK, posts)
}
