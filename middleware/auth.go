package middleware

import (
    "PA/utils"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

// AuthMiddleware godoc
// @Summary Middleware untuk autentikasi JWT
// @Description Memverifikasi token JWT dan melindungi endpoint yang memerlukan autentikasi
// @Tags Auth
// @Security ApiKeyAuth
// @Param Authorization header string true "Token JWT (Format: Bearer <token>)"
// @Success 200 {object} map[string]interface{} "Token valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ada"
// @Failure 500 {object} map[string]string "Kesalahan server internal"
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Token is required"})
            c.Abort()
            return
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")        

        token, user, err := utils.ParseJWT(tokenString)
        if err != nil || token == nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", user.ID)
        c.Next()
    }
}