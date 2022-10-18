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
// @Success 200 {integer} integer 1
// @Param input body services.CreatePostInput true "post information"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Failure default {object} errors.ErrorResponse
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
    postId, err := h.postService.Create(i)
    if err != nil {
        ResponseWithError(c, errors.NewHTTPError(http.StatusInternalServerError, err))
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id": postId,
    })

}

// @Summary getPostByid
// @Tags post
// @Description get one post by id
// @ID get-post-by-id
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Param        id   path      int  true  "Post ID"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Failure default {object} errors.ErrorResponse
// @Router /post/{id} [get]
func (h *PostHandler) GetOneById(c *gin.Context) {
    id := c.Param("id")
    postId, err := strconv.Atoi(id)
    if err != nil {
        ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
        return
    }

    post, err := h.postService.GetOneById(postId)
    if err != nil {
        ResponseWithError(c, errors.NewHTTPError(http.StatusInternalServerError, err))
        return
    }

    c.JSON(http.StatusOK, post)
}

// @Summary GetAllUserPosts
// @Tags post
// @Description Get All User's Posts
// @ID get-all-user-posts
// @Accept json
// @Produce json
// @Success 200 {array} []repository.Post
// @Param username path string true "User username"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Failure default {object} errors.ErrorResponse
// @Router /user/{username}/posts [get]
func (h *PostHandler) GetAllUserPosts(c *gin.Context) {
    username := c.Param("username")

    posts, err := h.postService.GetAllUserPosts(username)
    if err != nil {
        ResponseWithError(c, errors.NewHTTPError(http.StatusInternalServerError, err))
        return
    }

    c.JSON(http.StatusOK, posts)
}

type DeletePostResponse struct {
    Ok bool `json:"ok"`
}

// @Summary DeleteById
// @Tags post
// @Description delete one post by id
// @Security ApiKeyAuth
// @ID delete-post-by-id
// @Accept json
// @Produce json
// @Success 200 {object} DeletePostResponse
// @Param        id   path      int  true  "Post ID"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Failure default {object} errors.ErrorResponse
// @Router /post/{id} [delete]
func (h *PostHandler) DeleteById(c *gin.Context) {
    userId, exist := c.Get(userCtx)
    if !exist {
        ResponseWithError(c, errors.ErrIdNotFound)
        return
    }

    id := c.Param("id")
    postId, err := strconv.Atoi(id)
    if err != nil {
        ResponseWithError(c, errors.ErrInvalidPostId)
        return
    }

    post, err := h.postService.GetOneById(postId)
    if err != nil {
        ResponseWithError(c, errors.ErrPostNotFound)
        return
    }

    if post.UserId != userId {
        ResponseWithError(c, errors.ErrUserIsNotAuthor)
        return
    }

    err = h.postService.DeleteById(postId)
    if err != nil {
        ResponseWithError(c, errors.NewHTTPError(http.StatusInternalServerError, err))
        return
    }

    c.JSON(http.StatusOK, DeletePostResponse{
        Ok: true,
    })
}
