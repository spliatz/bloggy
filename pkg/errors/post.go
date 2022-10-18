package errors

import (
    "errors"
    "net/http"
)

var (
    ErrPostNotFound = NewHTTPError(http.StatusNotFound, errors.New("post does not exist"))

    ErrUserIsNotAuthor = NewHTTPError(http.StatusForbidden, errors.New("the user is not the author of the post"))

    ErrEmptyContent  = NewHTTPError(http.StatusBadRequest, errors.New("content must be not empty"))
    ErrInvalidPostId = NewHTTPError(http.StatusBadRequest, errors.New("invalid post id"))
)
