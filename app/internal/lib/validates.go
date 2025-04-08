package lib

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateLogin(fl validator.FieldLevel) bool {
	login := fl.Field().String()
	if isEmail := ValidateEmail(login); isEmail {
		return true
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

	return len(login) >= 5 && len(login) <= 15 && usernameRegex.MatchString(login)
}

func ValidateEmail(value string) bool {
	emailRegx := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if emailRegx.MatchString(value) {
		return true
	}
	return false
}
