package utils

import (
	"time"
	"os"
	"log"
	"errors"
	"PA/models"
	
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error tidak dapat membaca .env")
	}
}

var SecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func ParseJWT(tokenString string) (*jwt.Token, *models.User, error) {
    token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return SecretKey, nil
    })

    if err != nil {
        return nil, nil, err
    }

    if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
        if id, ok := (*claims)["id"].(float64); ok {
            user := &models.User{
                ID: uint(id),
                Username: (*claims)["username"].(string),
            }
            return token, user, nil
        }
    }

    return nil, nil, errors.New("Invalid token claims")
}