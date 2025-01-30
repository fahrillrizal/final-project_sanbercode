package main

import (
	"log"
	"os"
	
	"github.com/joho/godotenv"
	"PA/database"
	"PA/routes"
)

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