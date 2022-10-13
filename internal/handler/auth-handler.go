package handler

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
    "github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
    e "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type authHandler struct {
    authService services.Auth
}

func newAuthHandler(s services.Auth) *authHandler {
    return &authHandler{
        authService: s,
    }
}

func (h *authHandler) signUp(c *gin.Context) {

    req := structs.SignUpRequest{}
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err,
        })

        return
    }

    id, err := h.authService.SignUp(&req)
    if errors.Is(err, e.ErrShortPass) || errors.Is(err, e.ErrSimplePass) || errors.Is(err, e.ErrWrongUsername) ||
        errors.Is(err, e.ErrWrongUsernameLength) || errors.Is(err, e.ErrTakenUsername) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })

        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })

        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id": id,
    })

}
