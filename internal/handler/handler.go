package handler

import (
    "github.com/gin-gonic/gin"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
    a "github.com/Intellect-Bloggy/bloggy-backend/pkg/auth"
)

type Handlers struct {
    user
    auth
}

type user interface {
    getUserByUsername(ctx *gin.Context)
    editUser(ctx *gin.Context)
}

type auth interface {
    signUp(ctx *gin.Context)
    signIn(ctx *gin.Context)
    refresh(ctx *gin.Context)

    // Middlewares
    userIdentity(ctx *gin.Context)
}

func NewHandlers(s *services.Services, t a.TokenManager) *Handlers {
    return &Handlers{
        user: newUserHandler(s.User),
        auth: newAuthHandler(s.Auth, t),
    }
}

func (h *Handlers) InitRoutes() *gin.Engine {
    router := gin.New()

    auth := router.Group("/auth")
    {
        auth.POST("/signup", h.signUp)
        auth.POST("/signin", h.signIn)
        auth.POST("/refresh", h.refresh)
    }

    return router
}
