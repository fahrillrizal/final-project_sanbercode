package repository

import (
	"PA/models"

	"gorm.io/gorm"
)

func GetUserByUsername(db *gorm.DB, username string) (models.User, error) {
    var user models.User
    err := db.Where("username = ?", username).First(&user).Error
    return user, err
}

func GetUserByEmail(db *gorm.DB, email string) (models.User, error) {
    var user models.User
    err := db.Where("email = ?", email).First(&user).Error
    return user, err
}

func CreateUser(db *gorm.DB, user models.User) error {
	return db.Create(&user).Error
}