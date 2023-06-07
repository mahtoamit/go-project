package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title    string `json:"title" validate:"required"`
	Author   string `json:"author" validate:"required"`
	Rating   int    `json:"rating" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	gorm.Model    
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Location string `json:"location" validate:"required"`
}
