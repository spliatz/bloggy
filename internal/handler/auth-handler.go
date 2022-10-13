package handler

import (
	"github.com/Intellect-Bloggy/bloggy-backend/internal/services"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/structs"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authHandler struct {
	authService services.Auth
}

func newAuthHandler(s services.Auth) *authHandler {
	return &authHandler{
		authService: s,
	}
}

func (h *authHandler) registration(c *gin.Context) {
	input := structs.UserCreateInput{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Некорректное тело запроса",
		})
		return
	}

	id, err := h.authService.Registration(&input)

	if err != nil {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": err,
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]int{
		"id": id,
	})

}
