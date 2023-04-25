package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/spliatz/bloggy-backend/internal/controller/http/response"
	"github.com/spliatz/bloggy-backend/internal/domain/entity"
	auth_dto "github.com/spliatz/bloggy-backend/internal/domain/usecase/auth/dto"
	user_dto "github.com/spliatz/bloggy-backend/internal/domain/usecase/user/dto"
	"github.com/spliatz/bloggy-backend/pkg/errors"
)

type authUsecase interface {
	SignUp(ctx context.Context, dto user_dto.CreateUserDTO) (entity.Auth, error)
	SignIn(ctx context.Context, dto user_dto.GetByCredentialsDTO) (entity.Auth, error)
	Refresh(ctx context.Context, dto auth_dto.RefreshDTO) (entity.Auth, error)
	Logout(ctx context.Context, dto auth_dto.LogoutDTO) error
	UpdatePassword(ctx context.Context, dto auth_dto.UpdatePasswordDTO) error
}

type authHandler struct {
	authUsecase
	authMiddleware
}

func NewAuthHandler(authUsecase authUsecase, authMiddleware authMiddleware) *authHandler {
	return &authHandler{authUsecase: authUsecase, authMiddleware: authMiddleware}
}

func (h *authHandler) Register(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.signUp)
		auth.POST("/signin", h.signIn)
		auth.POST("/refresh", h.refresh)
		auth.DELETE("/logout", h.logout)
		protected := auth.Group("/account", h.authMiddleware.UserIdentity)
		{
			protected.PUT("/password", h.updatePassword)
		}
	}
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body dto.CreateUserDTO true "account info"
// @Success 201 {object} response.TokenResponse
// @Failure 400,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/signup [post]
func (h *authHandler) signUp(c *gin.Context) {
	dto := user_dto.CreateUserDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	res, err := h.authUsecase.SignUp(c.Request.Context(), dto)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.SetCookie(refreshTokenCookiesName, res.Refresh, monthInSeconds, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, response.TokenResponse{Access: res.Access})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body dto.GetByCredentialsDTO true "account username and password"
// @Success 200 {object} response.TokenResponse
// @Failure 400,403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/signin [post]
func (h *authHandler) signIn(c *gin.Context) {
	dto := user_dto.GetByCredentialsDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
		return
	}

	res, err := h.authUsecase.SignIn(c.Request.Context(), dto)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.SetCookie(refreshTokenCookiesName, res.Refresh, monthInSeconds, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, response.TokenResponse{Access: res.Access})
}

// @Summary Refresh
// @Tags auth
// @Description get new access and refresh token
// @ID get-new-access-and-refresh-token
// @Accept json
// @Produce json
// @Success 200 {object} response.TokenResponse
// @Failure 400,403,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *authHandler) refresh(c *gin.Context) {
	token, err := c.Cookie(refreshTokenCookiesName)
	if err != nil {
		response.ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
	}
	
	dto := auth_dto.RefreshDTO{RefreshToken: token}

	res, err := h.authUsecase.Refresh(c.Request.Context(), dto)
	if err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.SetCookie(refreshTokenCookiesName, res.Refresh, monthInSeconds, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, response.TokenResponse{Access: res.Access})
}

// @Summary Logout
// @Tags auth
// @Description logout user session
// @ID logout
// @Accept json
// @Produce json
// @Success 200
// @Failure 400,403,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/logout [delete]
func (h *authHandler) logout(c *gin.Context) {
	token, err := c.Cookie(refreshTokenCookiesName)
	if err != nil {
		response.ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
	}

	dto := auth_dto.LogoutDTO{RefreshToken: token}

	if err := h.authUsecase.Logout(c.Request.Context(), dto); err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.SetCookie(refreshTokenCookiesName, "", -1, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}

// @Summary UpdatePassword
// @Tags auth
// @Security ApiKeyAuth
// @Description update account password
// @ID update-password
// @Accept json
// @Produce json
// @Param input body dto.UpdatePasswordDTO true "old and new passwords"
// @Success 204
// @Failure 400,403,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/account/password [put]
func (h *authHandler) updatePassword(c *gin.Context) {
	userIdI, exist := c.Get(fieldUserId)
	if !exist {
		response.ResponseWithError(c, errors.ErrIdNotFound)
		return
	}

	userId, _ := userIdI.(int)

	dto := auth_dto.UpdatePasswordDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	dto.UserId = userId

	if err := h.UpdatePassword(c.Request.Context(), dto); err != nil {
		response.ResponseWithError(c, errors.EtoHe(err))
		return
	}

	c.Status(http.StatusNoContent)
}
