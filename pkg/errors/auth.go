package errors

import (
    "errors"
)

var ErrWrongPassOrUsername = errors.New("wrong username or password")
var ErrShortPass = errors.New("password must be at least 8 characters")
var ErrSimplePass = errors.New("password must contain a mix of letters, numbers, and symbols")
var ErrTokenExpired = errors.New("token expired")
var WrongToken = errors.New("token not found")
