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

func ResponseWithError(c *gin.Context, err errors.HTTPError) {
    logrus.Error(fmt.Sprintf(`[%d] %s`, err.Status(), err.Error()))

    c.AbortWithStatusJSON(err.Status(), ErrorResponse{
        Error: err.Error(),
    })
}
