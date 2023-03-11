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
	GetByUsername(ctx context.Context, username string) (entity.UserResponse, error)
	EditById(ctx context.Context, id int, dto user_dto.EditUserDTO) (entity.UserResponse, error)
	GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error)
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
	user, err := h.userUsecase.GetByUsername(c.Request.Context(), username)
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
	posts, err := h.userUsecase.GetAllByUsername(c.Request.Context(), username)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.JSON(http.StatusOK, posts)
}
