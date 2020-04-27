package auth

import (
	"unicode"

	cerror "github.com/x-color/calendar/model/error"
	"golang.org/x/crypto/bcrypt"
)

func validateSigninInfo(name, password string) error {
	if name == "" {
		return cerror.NewInvalidContentError(
			nil,
			"name is empty",
		)
	}
	if !isValidPassword(password) {
		return cerror.NewInvalidContentError(
			nil,
			"invalid password",
		)
	}
	return nil
}

func isValidPassword(password string) bool {
	hasMinLen := false
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false
	if 7 < len(password) && len(password) < 73 {
		hasMinLen = true
	}
	for _, char := range password {
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

func passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func verifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
