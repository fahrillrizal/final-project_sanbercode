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

func RegisterService(db *gorm.DB, user models.User) error {
	if user.Email != "" && !utils.IsValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	if strings.Contains(user.Username, "@") {
        return errors.New("username cannot contain '@'")
    }

	if user.Password != "" && !utils.IsValidPassword(user.Password) {
		return errors.New("invalid password format")
	}

	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("gagal Hash Password")
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