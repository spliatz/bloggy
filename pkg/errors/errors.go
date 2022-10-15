package errors

import (
    "errors"
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
