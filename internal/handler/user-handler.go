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

// @Summary GetUserByUsername
// @Tags user
// @Description login
// @ID get-user-by-username
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param username path string true "User username"
// @Success 200 {object} services.UserResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Failure default {object} errors.ErrorResponse
// @Router /user/{username} [get]
func (h *UserHandler) getUserByUsername(c *gin.Context) {
    username := c.Param("username")

    user, err := h.userService.GetUserByUsername(c.Request.Context(), username)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
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
// @Success 200 {object} services.UserResponse
// @Param input body services.EditInput true "user information"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Failure default {object} errors.ErrorResponse
// @Router /user [patch]
func (h *UserHandler) editUser(c *gin.Context) {
    c.Set("user_id", 1)
    userIdI, exist := c.Get("user_id")
    if !exist {
        ResponseWithError(c, errors.ErrIdNotFound)
        return
    }
    userId, _ := userIdI.(int)

    i := services.EditInput{}
    if err := c.BindJSON(&i); err != nil {
        ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
        return
    }

    user, err := h.userService.EditById(c.Request.Context(), userId, i)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusOK, user)
}
