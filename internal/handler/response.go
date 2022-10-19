package handler

import (
    "fmt"

    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

type ErrorResponse struct {
    Error string `json:"error"`
}

type EmptyResponse struct {
}

type IdResponse struct {
    Id int `json:"id"`
}

type TokenResponse struct {
    Access  string `json:"access_token"`
    Refresh string `json:"refresh_token"`
}

func ResponseWithError(c *gin.Context, err errors.HTTPError) {
    logrus.Error(fmt.Sprintf(`[%d] %s`, err.Status(), err.Error()))

    c.AbortWithStatusJSON(err.Status(), ErrorResponse{
        Error: err.Error(),
    })
}
