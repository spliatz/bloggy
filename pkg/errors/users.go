package errors

import (
    "errors"
)

var ErrWrongUsername = errors.New("username must contain only letters, numbers and some special characters")
var ErrWrongUsernameLength = errors.New("username must be 3 to 30 characters long")
var ErrUsernameNotFound = errors.New("username does not exist")
var ErrWrongId = errors.New("id must be non negative integer")
var ErrIdNotFound = errors.New("id does not exist")
var ErrUserDoesNotExist = errors.New("user does not exist")
var ErrTakenUsername = errors.New("a user with this username already exists")
var ErrTakenEmail = errors.New("a user with this email already exists")
var ErrTakenPhone = errors.New("a user with this phone already exists")
var ErrWrongNameLength = errors.New("name must be no longer than 60 characters")
var ErrWrongName = errors.New("name must contain only letters and spaces")
var ErrWrongDateFormat = errors.New(`date must be "2000-12-31" format`)
