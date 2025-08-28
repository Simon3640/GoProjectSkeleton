package models

import (
	"regexp"
	"unicode"
)

func IsValidPassword(p string) bool {
	var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool
	if len(p) >= 8 {
		hasMinLen = true
	}

	for _, char := range p {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func IsValidEmail(email string) bool {
	// regex for email validation
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email)
}

type HasValidation interface {
	Validate() []string
}

type HasUserID interface {
	GetUserID() uint
}
