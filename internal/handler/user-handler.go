package handler

import (
    "net/http"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService services.User
}

func newUserHandler(userService services.User) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

func (h *UserHandler) getUserByUsername(c *gin.Context) {
    username := c.Param("username")

    user, err := h.userService.GetUserByUsername(c.Request.Context(), username)
    if errors.Is(err, errors.ErrWrongUsername) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    if errors.Is(err, errors.ErrUsernameNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
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

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) editUser(c *gin.Context) {
    c.Set("user_id", 1)
    userIdI, exist := c.Get("user_id")
    if !exist {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": errors.ErrWrongId.Error(),
        })
        return
    }
    userId, _ := userIdI.(int)

    eReq := services.EditInput{}
    if err := c.BindJSON(&eReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    user, err := h.userService.EditById(c.Request.Context(), userId, eReq)
    if errors.IsOneOf(
        err,
        errors.ErrTakenUsername, errors.ErrTakenEmail, errors.ErrTakenPhone,
        errors.ErrWrongNameLength, errors.ErrWrongName,
    ) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    if errors.Is(err, errors.ErrIdNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
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

    c.JSON(http.StatusOK, user)
}
