package handler

import (
	"github.com/Intellect-Bloggy/bloggy-backend/internal/services"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
	"github.com/gin-gonic/gin"
	"net/http"
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
	user := structs.User{}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Некорректное тело запроса",
		})
		return
	}

	newUser, err := h.userService.Create(user)

	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}
