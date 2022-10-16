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
        e.NewHTTPError(c, http.StatusNotFound, e.ErrUserDoesNotExist)
        return
    }

    var createRequest services.CreatePostInput
    err := c.BindJSON(&createRequest)
    if err != nil {
        e.NewHTTPError(c, http.StatusBadRequest, e.ErrContentNotFound)
        return
    }

    createRequest.UserId = userId.(int)
    postId, err := h.postService.Create(createRequest)
    if err != nil {
        e.NewHTTPError(c, http.StatusInternalServerError, err)
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
        e.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }

    post, err := h.postService.GetOne(postId)
    if err != nil {
        e.NewHTTPError(c, http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetAllUserPosts(c *gin.Context) {
    username := c.Param("username")

    posts, err := h.postService.GetAllUserPosts(username)
    if err != nil {
        e.NewHTTPError(c, http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) Delete(c *gin.Context) {
    userId, exist := c.Get(userCtx)
    if !exist {
        e.NewHTTPError(c, http.StatusBadRequest, e.ErrUserDoesNotExist)
        return
    }

    id := c.Param("id")
    postId, err := strconv.Atoi(id)
    if err != nil {
        e.NewHTTPError(c, http.StatusBadRequest, e.ErrInvalidPostId)
        return
    }

    post, err := h.postService.GetOne(postId)
    if err != nil {
        e.NewHTTPError(c, http.StatusBadRequest, e.ErrPostNotFound)
        return
    }

    if post.UserId != userId {
        e.NewHTTPError(c, http.StatusUnauthorized, e.ErrUserIsNotAuthor)
        return
    }

    err = h.postService.Delete(postId)
    if err != nil {
        e.NewHTTPError(c, http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "ok": true,
    })
}
