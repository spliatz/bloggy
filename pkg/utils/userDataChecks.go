package utils

import (
	"github.com/spliatz/bloggy-backend/pkg/errors"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const emailCheckRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
const phoneCheckRegex = `^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`

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

func CheckName(s string) error {
	if len(s) > 60 {
		return errors.ErrWrongNameLength
	}

	for _, c := range []rune(s) {
		if !unicode.IsLetter(c) && !(c == ' ') {
			return errors.ErrWrongName
		}
	}

	return nil
}

func CheckEmail(email string) error {
	emailRegex := regexp.MustCompile(emailCheckRegex)
	if ok := emailRegex.MatchString(email); !ok || len(email) > 50 {
		return errors.ErrWrongEmail
	}
	return nil
}

func CheckPhone(phone string) error {
	emailRegex := regexp.MustCompile(phoneCheckRegex)
	if ok := emailRegex.MatchString(phone); !ok || len(phone) > 15 {
		return errors.ErrWrongPhone
	}
	return nil
}

func ParseDate(s string) (time.Time, error) {
	const dateFormat = "2006-01-02"
	date, err := time.Parse(dateFormat, s)
	if err != nil {
		return time.Time{}, errors.ErrWrongDateFormat
	}
	return date, nil
}

func IsLatin(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
