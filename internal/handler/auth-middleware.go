package handler

import (
    "errors"
    "net/http"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
)

const (
    authorizationHeader = "Authorization"

    userCtx = "user_id"
)

func (h *authHandler) userIdentity(c *gin.Context) {
    id, err := h.parseAuthHeader(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": err.Error(),
        })

        return
    }

    c.Set(userCtx, id)
}

func (h *authHandler) parseAuthHeader(c *gin.Context) (int, error) {
    header := c.GetHeader(authorizationHeader)
    if header == "" {
        return 0, errors.New("empty auth header")
    }

    headerParts := strings.Split(header, " ")
    if len(headerParts) != 2 || headerParts[0] != "Bearer" {
        return 0, errors.New("invalid auth header")
    }

    if len(headerParts[1]) == 0 {
        return 0, errors.New("token is empty")
    }

    idS, err := h.tokenManager.Parse(headerParts[1])
    if err != nil {
        return 0, err
    }

    id, err := strconv.Atoi(idS)
    if err != nil {
        return 0, err
    }

    return id, nil
}
