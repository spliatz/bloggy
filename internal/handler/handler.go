package handler

import (
    "github.com/gin-gonic/gin"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
)

type Handlers struct {
    user
    auth
}

type user interface {
    create(ctx *gin.Context)
}

type auth interface {
    registration(ctx *gin.Context)
}

func NewHandlers(s *services.Services) *Handlers {
    return &Handlers{
        user: newUserHandler(s.User),
        auth: newAuthHandler(s.Auth),
    }
}

func (h *Handlers) InitRoutes() *gin.Engine {
    router := gin.New()

    router.POST("/registration", h.registration)

    return router
}
