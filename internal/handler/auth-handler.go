package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
    a "github.com/Intellect-Bloggy/bloggy-backend/pkg/auth"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type authHandler struct {
    authService  services.Auth
    tokenManager a.TokenManager
}

func newAuthHandler(s services.Auth, t a.TokenManager) *authHandler {
    return &authHandler{
        authService:  s,
        tokenManager: t,
    }
}

type tokenResponse struct {
    Access  string `json:"access_token"`
    Refresh string `json:"refresh_token"`
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body services.SignUpInput true "account info"
// @Success 201 {object} tokenResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Failure default {object} errors.ErrorResponse
// @Router /auth/signup [post]
func (h *authHandler) signUp(c *gin.Context) {
    var i services.SignUpInput
    if err := c.BindJSON(&i); err != nil {
        ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
        return
    }

    t, err := h.authService.SignUp(c.Request.Context(), i)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusCreated, tokenResponse{t.Access, t.Refresh})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body services.SignInInput true "account username and password"
// @Success 201 {object} tokenResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Failure default {object} errors.ErrorResponse
// @Router /auth/signin [post]
func (h *authHandler) signIn(c *gin.Context) {
    var i services.SignInInput
    if err := c.BindJSON(&i); err != nil {
        ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
        return
    }

    res, err := h.authService.SignIn(c.Request.Context(), i)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusOK, tokenResponse{res.Access, res.Refresh})
}

type refreshInput struct {
    Token string `json:"token" binding:"required"`
}

// @Summary Refresh
// @Tags auth
// @Description get new access and refresh token
// @ID get-new-access-and-refresh-token
// @Accept json
// @Produce json
// @Param input body refreshInput true "refresh token"
// @Success 201 {object} tokenResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Failure default {object} errors.ErrorResponse
// @Router /auth/refresh [post]
func (h *authHandler) refresh(c *gin.Context) {
    var i refreshInput
    if err := c.BindJSON(&i); err != nil {
        ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
        return
    }

    res, err := h.authService.RefreshTokens(c.Request.Context(), i.Token)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusOK, tokenResponse{
        Access:  res.Access,
        Refresh: res.Refresh,
    })
}
