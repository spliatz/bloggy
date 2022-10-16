package errors

import (
    "errors"
)

var ErrContentNotFound = errors.New("content is required option")
var ErrPostNotFound = errors.New("post does not exist")
var ErrUserIsNotAuthor = errors.New("post can be deleted only by author")
var ErrInvalidPostId = errors.New("invalid post id")
