package main

import (
	"PA/controllers"
	"PA/middleware"
	"PA/models"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

var (
	db *gorm.DB
	err error
)
func main() {
	err = godotenv.Load(".env")
    if err != nil {
        panic("Error loading .env file")
    }

	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
        os.Getenv("PGHOST"),
        os.Getenv("PGPORT"),
        os.Getenv("PGUSER"),
        os.Getenv("PGPASSWORD"),
        os.Getenv("PGDATABASE"),
    )

	db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if err != nil {
        panic(err)
    }

	err = db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.ProjectCollaborator{},
		&models.Task{},
		&models.TaskAssignment{},
	)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Use(func (c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	projects := router.Group("/projects")
	projects.Use(middleware.AuthMiddleware())
	{
		projects.POST("/", controllers.AddProjectController)
		projects.GET("/", controllers.GetProjectsController)
		projects.GET("/:id", controllers.GetProjectByIDController)
		projects.PUT("/:id", controllers.EditProjectController)
		projects.DELETE("/:id", controllers.DeleteProjectController)
		projects.POST("/:id/collaborators", controllers.AddCollaboratorController)
		projects.DELETE("/:id/collaborators", controllers.RemoveCollaboratorController)
	}


	router.Run(":8080")
}