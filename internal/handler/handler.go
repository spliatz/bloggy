package handler

import (
	"github.com/Intellect-Bloggy/bloggy-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	user
}

type user interface {
	create(ctx *gin.Context)
}

func NewHandlers(s *services.Services) *Handlers {
	return &Handlers{
		user: newUserHandler(s.User),
	}
}

func (h *Handlers) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("api")
	{
		users := api.Group("users")
		{
			users.POST("/create", h.user.create)
		}
	}

	return router
}
