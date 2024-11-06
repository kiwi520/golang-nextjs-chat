package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(fromPassword), nil
}

func VerifyPassword(password string, hashedPassword string) error {
	fmt.Println("password")
	fmt.Println(password)
	fmt.Println("password")
	fmt.Println("hashedPassword")
	fmt.Println(hashedPassword)
	fmt.Println("hashedPassword")
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
