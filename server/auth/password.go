package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength int = 8
	bCryptCost        int = 10

	// if the user's password hash is empty, use this
	// hash to mask the fact
	timingHash = "$2a$10$4Kys.PIxpCIoUmlcY6D7QOTuMPgk27lpmV74OWCWfqjwnG/JN4kcu"
)

var (
	ErrPasswordResetExpired = errors.New("password reset link is expired")
)

// bcrypt compare hash with given password
func VerifyPassword(hash, password string) bool {
	// incase either hash or password is empty, compare
	// something and return false to mask the timing
	if len(hash) == 0 || len(password) == 0 {
		if err := bcrypt.CompareHashAndPassword([]byte(timingHash), []byte(password)); err != nil {
			return false
		}
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func encrypt(password string) ([]byte, error) {
	// encrypt with bcrypt
	return bcrypt.GenerateFromPassword([]byte(password), bCryptCost)
}
