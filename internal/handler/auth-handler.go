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

func (h *authHandler) signUp(c *gin.Context) {
    var i services.SignUpInput
    if err := c.BindJSON(&i); err != nil {
        errors.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }

    t, err := h.authService.SignUp(c.Request.Context(), i)
    if err != nil {
        if errors.IsOneOf(
            err,
            errors.ErrShortPass, errors.ErrSimplePass, errors.ErrWrongUsername,
            errors.ErrWrongUsernameLength, errors.ErrTakenUsername, errors.ErrWrongDateFormat,
        ) {
            errors.NewHTTPError(c, http.StatusBadRequest, err)
            return
        }

        errors.NewHTTPError(c, http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusCreated, tokenResponse{t.Access, t.Refresh})
}

func (h *authHandler) signIn(c *gin.Context) {
    var i services.SignInInput
    if err := c.BindJSON(&i); err != nil {
        errors.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }

    res, err := h.authService.SignIn(c.Request.Context(), i)
    if errors.Is(err, errors.ErrWrongPassOrUsername) {
        errors.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }
    if errors.Is(err, errors.ErrTokenExpired) {
        errors.NewHTTPError(c, http.StatusUnauthorized, err)
        return
    }
    if err != nil {
        errors.NewHTTPError(c, http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, tokenResponse{res.Access, res.Refresh})
}

type refreshInput struct {
    Token string `json:"token" binding:"required"`
}

func (h *authHandler) refresh(c *gin.Context) {
    var i refreshInput
    if err := c.BindJSON(&i); err != nil {
        errors.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }

    res, err := h.authService.RefreshTokens(c.Request.Context(), i.Token)
    if errors.Is(err, errors.WrongToken) {
        errors.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }
    if err != nil {
        errors.NewHTTPError(c, http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, tokenResponse{
        Access:  res.Access,
        Refresh: res.Refresh,
    })
}
