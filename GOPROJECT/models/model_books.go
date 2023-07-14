package models

import (
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title     string `json:"title" validate:"required"`
	Author    string `json:"author" validate:"required"`
	Rating    int    `json:"rating" validate:"required"`
	OtherInfo string `json:"otherinfo" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required" `
}

type SignupRequest struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Location string `json:"location" validate:"required"`
}

// Custom validation tag for email
func ValidateEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

// Custom validation tag for password
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasSpecialChar := false
	hasNumber := false
	hasUpper := false
	hasLower := false

	for _, char := range password {
		switch {
		case strings.ContainsRune("!@#$%^&*()-_+=~`{[}]|:;<>,.?/", char):
			hasSpecialChar = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		}
	}

	return hasSpecialChar && hasNumber && hasUpper && hasLower
}

type CacheModel struct {
	DatabaseName string `validate:"required"`
	QueryType    string `validate:"required"`
	Model        Book
}
