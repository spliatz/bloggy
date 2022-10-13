package errors

import (
    "errors"
)

var ErrWrongUsername = errors.New("username must contain only letters, numbers and some special characters")
var ErrWrongUsernameLength = errors.New("username must be 3 to 30 characters long")
var ErrUsernameNotFound = errors.New("username does not exist")
var ErrTakenUsername = errors.New("a user with this name already exists")
