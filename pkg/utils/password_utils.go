package utils

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
func ComparePasswords(hashedPwd []byte, plainPwd []byte) bool {
	if len(hashedPwd) == 0 || len(plainPwd) == 0 {
		return false
	}
	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	return err == nil
}

func SaltAndVerifyPassword(s string) (pw []byte, err error) {
	if s == "" {
		return nil, errors.New("password cannot be empty")
	}

	letters := 0
	number := false
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
		}
	}
	if letters <= 4 {
		return nil, errors.New("password must contain at least 5 letters")
	}
	if !number {
		return nil, errors.New("password must contain at least 1 number")
	}
	length := len(s)
	if length < 8 || length > 30 {
		return nil, errors.New("password must be between 8 and 30 characters")
	}

	hashedPw, err := HashAndSalt([]byte(s))
	if err != nil {
		return nil, err
	}
	return hashedPw, nil
}
