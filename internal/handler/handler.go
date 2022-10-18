package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"

    _ "github.com/Intellect-Bloggy/bloggy-backend/docs"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
    a "github.com/Intellect-Bloggy/bloggy-backend/pkg/auth"
)

type Handlers struct {
    user
    auth
    post
}

type user interface {
    GetByUsername(ctx *gin.Context)
    EditById(ctx *gin.Context)
}

type auth interface {
    SignUp(ctx *gin.Context)
    SignIn(ctx *gin.Context)
    Refresh(ctx *gin.Context)

    // Middlewares
    UserIdentity(ctx *gin.Context)
}

type post interface {
    Create(c *gin.Context)
    GetById(c *gin.Context)
    GetAllByUsername(c *gin.Context)
    DeleteById(c *gin.Context)
}

func NewHandlers(s *services.Services, t a.TokenManager) *Handlers {
    return &Handlers{
        user: newUserHandler(s.User),
        auth: newAuthHandler(s.Auth, t),
        post: newPostHandler(s.Post),
    }
}

func (h *Handlers) InitRoutes() *gin.Engine {
    router := gin.New()

    auth := router.Group("/auth")
    {
        auth.POST("/signup", h.auth.SignUp)
        auth.POST("/signin", h.auth.SignIn)
        auth.POST("/refresh", h.auth.Refresh)
    }

    user := router.Group("/user")
    {
        authProtectedUsers := user.Group("", h.auth.UserIdentity)
        {
            authProtectedUsers.PATCH("", h.user.EditById)
        }

        user.GET("/:username", h.user.GetByUsername)
        user.GET("/:username/posts", h.post.GetAllByUsername)
    }

    post := router.Group("/post")
    {
        authProtectedPost := post.Group("", h.auth.UserIdentity)
        {
            authProtectedPost.POST("", h.post.Create)
            authProtectedPost.DELETE("/:id", h.post.DeleteById)
        }

        post.GET("/:id", h.post.GetById)
    }

    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return router
}
