package errors

import (
    "errors"

    "github.com/gin-gonic/gin"
)

func Is(err, target error) bool {
    return errors.Is(err, target)
}

func IsOneOf(err error, targets ...error) bool {
    for _, target := range targets {
        if errors.Is(err, target) {
            return true
        }
    }

    return false
}

func NewHTTPError(c *gin.Context, status int, err error) {
    c.AbortWithStatusJSON(status, gin.H{
        "error": err.Error(),
    })
}
