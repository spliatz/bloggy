package handler

import (
    "github.com/gin-gonic/gin"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
    a "github.com/Intellect-Bloggy/bloggy-backend/pkg/auth"
)

type Handlers struct {
    user
    auth
    post
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

type post interface {
    Create(c *gin.Context)
    GetOne(c *gin.Context)
    GetAllUserPosts(c *gin.Context)
    Delete(c *gin.Context)
}

func NewHandlers(s *services.Services, t a.TokenManager) *Handlers {
    return &Handlers{
        user: newUserHandler(s.User),
        auth: newAuthHandler(s.Auth, t),
        post: newPostHandler(s.Post, s.User),
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

    user := router.Group("/user")
    {
        authProtectedUsers := user.Group("", h.userIdentity)
        {
            authProtectedUsers.PATCH("", h.editUser)
            authProtectedUsers.GET("/:username", h.getUserByUsername)
        }

        user.GET("/:username/posts", h.post.GetAllUserPosts)
    }

    post := router.Group("/post")
    {
        authProtectedPost := post.Group("", h.userIdentity)
        {
            authProtectedPost.POST("", h.post.Create)
            authProtectedPost.DELETE("/:id", h.post.Delete)
        }

        post.GET("/:id", h.post.GetOne)
    }

    return router
}
