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

// Register godoc
// @Summary Register new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.UserAuth true "Registration Data"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/register [post]
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

// Login godoc
// @Summary User login
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.UserAuth true "Login"
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "Bad Request - User not found or invalid password"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/login [post]
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

