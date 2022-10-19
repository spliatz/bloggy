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

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body services.SignUpInput true "account info"
// @Success 201 {object} TokenResponse
// @Failure 400,409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /auth/signup [post]
func (h *authHandler) SignUp(c *gin.Context) {
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

    c.JSON(http.StatusCreated, TokenResponse{t.Access, t.Refresh})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body services.SignInInput true "account username and password"
// @Success 200 {object} TokenResponse
// @Failure 400,403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /auth/signin [post]
func (h *authHandler) SignIn(c *gin.Context) {
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

    c.JSON(http.StatusOK, TokenResponse{res.Access, res.Refresh})
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
// @Success 200 {object} TokenResponse
// @Failure 400,403,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *authHandler) Refresh(c *gin.Context) {
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

    c.JSON(http.StatusOK, TokenResponse{
        Access:  res.Access,
        Refresh: res.Refresh,
    })
}
