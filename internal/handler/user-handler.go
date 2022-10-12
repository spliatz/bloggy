package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Intellect-Bloggy/bloggy-backend/internal/services"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
)

type UserHandler struct {
	userService services.User
}

func newUserHandler(userService services.User) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) create(c *gin.Context) {
	input := structs.UserCreateInput{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Некорректное тело запроса",
		})
		return
	}

	user, err := h.userService.Create(&input)

	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}
