package hasher

// Package hasher implements utility for
// hashing passwords

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns bcrypt hash of the password or error if password hash goes wrong.
// Use CheckPasswordHash, as defined in this package,
// to compare the returned hashed password with its cleartext version.
func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", fmt.Errorf("String must not be empty")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPasswordHash compares a bcrypt hashed password with its possible
// plaintext equivalent. Returns nil on success, or an error on failure.
func CheckPasswordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
