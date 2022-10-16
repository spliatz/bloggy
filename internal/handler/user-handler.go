package handler

import (
    "net/http"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
    e "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"

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
    if e.Is(err, e.ErrWrongUsername) {
        e.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }
    if e.Is(err, e.ErrUsernameNotFound) {
        e.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }
    if err != nil {
        e.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) editUser(c *gin.Context) {
    c.Set("user_id", 1)
    userIdI, exist := c.Get("user_id")
    if !exist {
        e.NewHTTPError(c, http.StatusBadRequest, e.ErrUserDoesNotExist)
        return
    }
    userId, _ := userIdI.(int)

    eReq := services.EditInput{}
    if err := c.BindJSON(&eReq); err != nil {
        e.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }

    user, err := h.userService.EditById(c.Request.Context(), userId, eReq)
    if e.IsOneOf(
        err,
        e.ErrTakenUsername, e.ErrTakenEmail, e.ErrTakenPhone,
        e.ErrWrongNameLength, e.ErrWrongName,
    ) {
        e.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }
    if e.Is(err, e.ErrIdNotFound) {
        e.NewHTTPError(c, http.StatusNotFound, err)
        return
    }
    if err != nil {
        e.NewHTTPError(c, http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, user)
}
