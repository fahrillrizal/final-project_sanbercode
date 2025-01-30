package routes

import (
	"PA/controllers"
	"PA/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	docs "PA/docs"
   	swaggerfiles "github.com/swaggo/files"
   	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.POST("/api/register", controllers.Register)
	router.POST("/api/login", controllers.Login)

	auth := router.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		setupProjectRoutes(auth)
		setupTaskRoutes(auth)
	}

	return router
}

func setupProjectRoutes(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	{
		projects.POST("/", controllers.AddProjectController)
		projects.GET("/", controllers.GetProjectsController)
		projects.GET("/:project_id", controllers.GetProjectByIDController)
		projects.PUT("/:project_id", controllers.EditProjectController)
		projects.DELETE("/:project_id", controllers.DeleteProjectController)
		
		projects.POST("/:project_id/collaborators", controllers.AddCollaboratorController)
		projects.DELETE("/:project_id/collaborators", controllers.RemoveCollaboratorController)

		tasks := projects.Group("/:project_id/tasks")
		{
			tasks.POST("/", controllers.AddTaskController)
			tasks.GET("/", controllers.GetTaskByProjectController)
			tasks.PUT("/:task_id", controllers.UpdateTaskController)
		}
	}
}

func setupTaskRoutes(rg *gin.RouterGroup) {
	rg.GET("/tasks", controllers.GetAllTaskController)
	rg.GET("/tasks/:id", controllers.GetTaskByIDController)
	rg.DELETE("/tasks/:id", controllers.DeleteTaskController)
}