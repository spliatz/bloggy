package errors

import (
    "errors"
    "net/http"
)

var (
    ErrTokenExpired  = NewHTTPError(http.StatusUnauthorized, errors.New("token expired"))
    ErrTokenNotFound = NewHTTPError(http.StatusUnauthorized, errors.New("token does not exist"))

    ErrWrongCredentials = NewHTTPError(http.StatusBadRequest, errors.New("wrong username or password"))
    ErrShortPass        = NewHTTPError(http.StatusBadRequest, errors.New("password must be at least 8 characters"))
    ErrSimplePass       = NewHTTPError(http.StatusBadRequest, errors.New("password must contain a mix of letters, numbers, and symbols"))
)
