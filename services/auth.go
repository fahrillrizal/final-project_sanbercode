package services

import (
	"errors"
	"PA/models"
	"PA/repository"
	"PA/utils"
	"strings"

	"gorm.io/gorm"
)

func LoginService(db *gorm.DB, identifier, password string) (string, error) {
	var user models.User
	var err error
	if strings.Contains(identifier, "@") {
		if utils.IsValidEmail(identifier) {
			user, err = repository.GetUserByEmail(db, identifier)
		} 
	} else {
		user, err = repository.GetUserByUsername(db, identifier)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user tidak ditemukan")
		}
		return "", err
	}

	if !utils.CheckPassword(password, user.Password) {
        return "", errors.New("invalid password")
    }

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return "", errors.New("gagal generate token")
	}

	return token, nil
}

func RegisterService(db *gorm.DB, input models.UserAuth) error {
    if input.Email != "" && !utils.IsValidEmail(input.Email) {
        return errors.New("invalid email format")
    }

    if strings.Contains(input.Username, "@") {
        return errors.New("username cannot contain '@'")
    }

    if !utils.IsValidPassword(input.Password) {
        return errors.New("invalid password format")
    }

    user := models.User{
        Username: input.Username,
        Email:    input.Email,
    }

    hashedPass, err := utils.HashPassword(input.Password)
    if err != nil {
        return errors.New("gagal hash password")
    }
    user.Password = hashedPass

    if err := repository.CreateUser(db, user); err != nil {
        return err
    }

    return nil
}

func ParseTokenService(tokenString string) (*models.User, error) {
	_, user, err := utils.ParseJWT(tokenString)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return user, nil
}