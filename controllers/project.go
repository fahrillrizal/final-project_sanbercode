package controllers

import (
	"net/http"
	"strconv"
	"PA/models"
	"PA/services"
	"strings"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProjectInput digunakan untuk validasi input add & edit project
type ProjectInput struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description"`
}


// Get Projects godoc
// @Summary Get all projects
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/projects [get]
func GetProjectsController(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)

	projects, err := services.GetAllProjectsService(db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": projects})
}

// Get Project by ID godoc
// @Summary Get project by ID
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param project_id path uint true "Project ID"
// @Success 200 {object} models.Project "Success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Project Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/projects/{project_id} [get]
func GetProjectByIDController(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	projectIDStr := c.Param("project_id")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := services.GetProjectByIDService(db, uint(projectID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})
}

// Add Project godoc
// @Summary Add a new project
// @Tags Projects
// @Security BearerAuth
// @Description Add a new project with a name, description, and owner
// @Accept json
// @Produce json
// @Param project body ProjectInput true "Project Input"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/projects [post]
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

// Edit Project godoc
// @Summary Edit a project
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param project_id path uint true "Project ID"
// @Param input body ProjectInput true "Project Data"
// @Success 200 {object} models.Project "Project Updated"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Project Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/projects/{project_id} [put]
func EditProjectController(c *gin.Context) {
	var input ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	projectIDStr := c.Param("project_id") 

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := services.GetProjectByIDService(db, uint(projectID), userID)
    if err != nil {
        errorMsg := err.Error()
        statusCode := http.StatusNotFound
        if errorMsg == "anda tidak memiliki akses ke project ini" {
            statusCode = http.StatusForbidden
        }
        c.JSON(statusCode, gin.H{"error": errorMsg})
        return
    }

	project.Name = input.Name
	project.Description = input.Description

	if err := services.UpdateProjectService(db, &project, userID); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{"data": project})
}

// Delete Project godoc
// @Summary Delete a project
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param project_id path uint true "Project ID"
// @Success 200 {object} map[string]string "Success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 403 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Project Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/projects/{project_id} [delete]
func DeleteProjectController(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    userID := c.MustGet("user_id").(uint)
    
    projectID, err := strconv.ParseUint(c.Param("project_id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }
    
    if err := services.DeleteProjectService(db, uint(projectID), userID); err != nil {
        errorMsg := gin.H{"error": err.Error()}
        status := http.StatusInternalServerError
        
        if strings.Contains(err.Error(), "unauthorized") {
            status = http.StatusForbidden
        } else if strings.Contains(err.Error(), "tidak ditemukan") {
            status = http.StatusNotFound
        }
        
        c.JSON(status, errorMsg)
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Project berhasil dihapus"})
}

// Add Collaborator godoc
// @Summary Add a collaborator to a project
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param project_id path uint true "Project ID"
// @Param input body CollaboratorInput true "Collaborator Data"
// @Success 200 {object} map[string]string "Collaborator added successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Project Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/projects/{project_id}/collaborators [post]
func AddCollaboratorController(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    ownerID := c.MustGet("user_id").(uint)

    projectIDStr := c.Param("project_id")
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

	// Memanggil service untuk add collaborator
    if err := services.AddCollaboratorService(db, uint(projectID), input.UserID, ownerID); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Collaborator berhasil ditambahkan"})
}

// Remove Collaborator godoc
// @Summary Remove a collaborator from a project
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param project_id path uint true "Project ID"
// @Param input body CollaboratorInput true "Collaborator Data"
// @Success 200 {object} map[string]string "Collaborator removed successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Project Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/projects/{project_id}/collaborators [delete]
func RemoveCollaboratorController(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    ownerID := c.MustGet("user_id").(uint)

    projectIDStr := c.Param("project_id")
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

	// Memanggil service untuk remove collaborator
    if err := services.RemoveCollaboratorService(db, uint(projectID), input.UserID, ownerID); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Collaborator berhasil dihapus"})
}