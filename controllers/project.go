package controllers

import (
	"net/http"
	"strconv"
	"PA/models"
	"PA/services"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectInput struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description"`
}


func GetProjectsController(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	projects, err := services.GetAllProjectsService(db)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": projects})
}

func GetProjectByIDController(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	projectIDStr := c.Param("id") 

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := services.GetProjectByIDService(db, uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if (models.Project{}) == project {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})
}

func AddProjectController(c *gin.Context) {
	var input ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)

	project := models.Project{
		Name: input.Name,
		Description: input.Description,
		OwnerID: userID,
	}

	if err := services.CreateProjectService(db, &project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": project})
}

func EditProjectController(c *gin.Context) {
	var input ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	projectIDStr := c.Param("id") 

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := services.GetProjectByIDService(db, uint(projectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project tidak ditemukan"})
		return
	}

	project.Name = input.Name
	project.Description = input.Description

	if err := services.UpdateProjectService(db, &project, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal untuk edit project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})
}

func DeleteProjectController(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	projectIDStr := c.Param("id") 

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := services.DeleteProjectService(db, uint(projectID), userID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project berhasil dihapus"})
}

func AddCollaboratorController(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	ownerID := c.MustGet("user_id").(uint)

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var input struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.AddCollaboratorService(db, uint(projectID), input.UserID, ownerID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collaborator berhasil ditambahkan"})
}

func RemoveCollaboratorController(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	ownerID := c.MustGet("user_id").(uint)

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var input struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.RemoveCollaboratorService(db, uint(projectID), input.UserID, ownerID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collaborator berhasil dihapus"})
}