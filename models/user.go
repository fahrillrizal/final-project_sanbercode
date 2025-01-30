package models

import "time"

type User struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Email string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type UserAuth struct {
    Username string `json:"username" binding:"required_without=Email"`
    Email    string `json:"email" binding:"required_without=Username"`
    Password string `json:"password" binding:"required"`
}