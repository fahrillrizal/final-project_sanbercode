package main

import (
	"log"
	"os"
	
	"github.com/joho/godotenv"
	"PA/database"
	"PA/routes"
)

// @title Project Management API
// @version 1.0
// @description API untuk manajemen proyek dan tugas
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	router := routes.SetupRouter(db)

	port := os.Getenv("PORT")
	
	router.Run(":" + port)
}