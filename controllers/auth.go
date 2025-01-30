package controllers

import (
	"PA/models"
	"PA/services"
	"PA/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
    var input models.UserAuth
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db := c.MustGet("db").(*gorm.DB)

    var identifier string
    if input.Username != "" {
        if strings.Contains(input.Username, "@") {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Username cannot contain '@'"})
            return
        }
        identifier = input.Username
    } else {
        if !utils.IsValidEmail(input.Email) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
            return
        }
        identifier = input.Email
    }

    token, err := services.LoginService(db, identifier, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

func Register(c *gin.Context) {
    var input models.UserAuth
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db := c.MustGet("db").(*gorm.DB)
    if err := services.RegisterService(db, input); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Registration successful. Please log in."})
}