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
        e.NewHTTPError(c, http.StatusBadRequest, err)
        return
    }

    post, err := h.postService.GetOneById(postId)
    if err != nil {
        e.NewHTTPError(c, http.StatusInternalServerError, err)
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
        e.NewHTTPError(c, http.StatusInternalServerError, err)
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
        e.NewHTTPError(c, http.StatusBadRequest, e.ErrUserDoesNotExist)
        return
    }

    id := c.Param("id")
    postId, err := strconv.Atoi(id)
    if err != nil {
        e.NewHTTPError(c, http.StatusBadRequest, e.ErrInvalidPostId)
        return
    }

    post, err := h.postService.GetOneById(postId)
    if err != nil {
        e.NewHTTPError(c, http.StatusBadRequest, e.ErrPostNotFound)
        return
    }

    if post.UserId != userId {
        e.NewHTTPError(c, http.StatusUnauthorized, e.ErrUserIsNotAuthor)
        return
    }

    err = h.postService.DeleteById(postId)
    if err != nil {
        e.NewHTTPError(c, http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, DeletePostResponse{
        Ok: true,
    })
}
