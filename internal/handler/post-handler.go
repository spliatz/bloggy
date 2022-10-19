package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type PostHandler struct {
    postService services.Post
}

func newPostHandler(sp services.Post) *PostHandler {
    return &PostHandler{
        postService: sp,
    }
}

// @Summary createPost
// @Tags post
// @Description create post
// @Security ApiKeyAuth
// @ID create-post
// @Accept json
// @Produce json
// @Param input body services.CreatePostInput true "post information"
// @Success 201 {object} IdResponse
// @Failure 400,401,403,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /post [post]
func (h *PostHandler) Create(c *gin.Context) {
    userId, exist := c.Get(userCtx)
    if !exist {
        ResponseWithError(c, errors.ErrIdNotFound)
        return
    }

    var i services.CreatePostInput
    err := c.BindJSON(&i)
    if err != nil {
        ResponseWithError(c, errors.ErrEmptyContent)
        return
    }

    i.AuthorId = userId.(int)
    postId, err := h.postService.Create(c.Request.Context(), i)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusCreated, IdResponse{postId})

}

// @Summary getPostByid
// @Tags post
// @Description get one post by id
// @ID get-post-by-id
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} services.PostResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /post/{id} [get]
func (h *PostHandler) GetById(c *gin.Context) {
    id := c.Param("id")
    postId, err := strconv.Atoi(id)
    if err != nil {
        ResponseWithError(c, errors.ErrWrongId)
        return
    }

    post, err := h.postService.GetById(c.Request.Context(), postId)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusOK, post)
}

// @Summary GetAllByUsername
// @Tags post
// @Description Get All User's Posts
// @ID get-all-by-username
// @Accept json
// @Produce json
// @Param username path string true "User username"
// @Success 200 {array} services.PostResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /user/{username}/posts [get]
func (h *PostHandler) GetAllByUsername(c *gin.Context) {
    username := c.Param("username")

    posts, err := h.postService.GetAllByUsername(c.Request.Context(), username)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusOK, posts)
}

// @Summary DeleteById
// @Tags post
// @Description delete one post by id
// @Security ApiKeyAuth
// @ID delete-post-by-id
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} EmptyResponse
// @Failure 400,401,403,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /post/{id} [delete]
func (h *PostHandler) DeleteById(c *gin.Context) {
    userIdI, exist := c.Get(userCtx)
    if !exist {
        ResponseWithError(c, errors.ErrIdNotFound)
        return
    }
    userId := userIdI.(int)

    id := c.Param("id")
    postId, err := strconv.Atoi(id)
    if err != nil {
        ResponseWithError(c, errors.ErrInvalidPostId)
        return
    }

    post, err := h.postService.GetById(c.Request.Context(), postId)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
        return
    }

    isAuthor, err := h.postService.IsAuthor(c.Request.Context(), post.Id, userId)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
        return
    }
    if !isAuthor {
        ResponseWithError(c, errors.ErrUserIsNotAuthor)
        return
    }

    err = h.postService.DeleteById(c.Request.Context(), postId)
    if err != nil {
        ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusOK, EmptyResponse{})
}
