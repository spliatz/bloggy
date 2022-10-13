package errors

import (
    "errors"
)

var ErrWrongPass = errors.New("wrong password")
var ErrShortPass = errors.New("password must be at least 8 characters")
var ErrSimplePass = errors.New("password must contain a mix of letters, numbers, and symbols")
