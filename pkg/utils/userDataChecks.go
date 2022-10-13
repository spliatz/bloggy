package utils

import (
    "strings"
    "unicode"

    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

func CheckUsername(s string) error {
    if len(s) < 3 || len(s) > 30 {
        return errors.ErrWrongUsernameLength
    }

    allowedChars := "_-."
    for _, c := range []rune(s) {
        if !unicode.IsDigit(c) &&
            !strings.ContainsRune(allowedChars, c) &&
            !IsLatin(c) {
            return errors.ErrWrongUsername
        }
    }

    return nil
}

func CheckPassword(s string) error {
    if len(s) < 8 {
        return errors.ErrShortPass
    }

    var contains_num, contains_letter, contains_sym bool

    allowedChars := "!@#$%^&*()-_=+{}[]\\|;:,<.>/?"
    for _, c := range []rune(s) {
        if unicode.IsDigit(c) {
            contains_num = true
        } else if strings.ContainsRune(allowedChars, c) {
            contains_sym = true
        } else if IsLatin(c) {
            contains_letter = true
        } else {
            // Сюда попадаем, когда появился запрещенный символ в пароле
            return errors.ErrSimplePass
        }
    }

    if !(contains_num && contains_letter && contains_sym) {
        // Сюда попадаем, если не соблюдены условия сложности пароля
        return errors.ErrSimplePass
    }

    return nil
}

func IsLatin(c rune) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
