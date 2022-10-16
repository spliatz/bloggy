package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
    e "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type PostHandler struct {
    postService services.Post
    userService services.User
}

func newPostHandler(sp services.Post, su services.User) *PostHandler {
    return &PostHandler{
        postService: sp,
        userService: su,
    }
}

func (h *PostHandler) Create(c *gin.Context) {
    userId, exist := c.Get(userCtx)
    if !exist {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": e.ErrUserDoesNotExist.Error(),
        })
        return
    }

    var createRequest services.CreatePostInput
    err := c.BindJSON(&createRequest)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": e.ErrContentNotFound.Error(),
        })
        return
    }

    createRequest.UserId = userId.(int)
    postId, err := h.postService.Create(createRequest)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id": postId,
    })

}

func (h *PostHandler) GetOne(c *gin.Context) {
    id := c.Param("id")
    postId, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid post id",
        })
        return
    }

    post, err := h.postService.GetOne(postId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": e.ErrPostNotFound.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetAllUserPosts(c *gin.Context) {
    username := c.Param("username")

    posts, err := h.postService.GetAllUserPosts(username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) Delete(c *gin.Context) {
    userId, exist := c.Get(userCtx)
    if !exist {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": e.ErrUserDoesNotExist.Error(),
        })
        return
    }

    id := c.Param("id")
    postId, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid post id",
        })
        return
    }

    post, err := h.postService.GetOne(postId)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": e.ErrPostNotFound.Error(),
        })
        return
    }

    if post.UserId != userId {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": e.ErrUserIsNotAuthor.Error(),
        })
        return
    }

    err = h.postService.Delete(postId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "ok": true,
    })
}
