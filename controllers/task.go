package controllers

import (
    "net/http"
    "strconv"
    "time"
    "PA/models"
    "PA/services"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// taskInput digunakan untuk validasi input add & edit task
type taskInput struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Status      string   `json:"status"`
    AssignedTo  []uint   `json:"assigned_to"`
    Deadline    string   `json:"deadline"`
}

// digunakan untuk parse input deadline
func parseDeadline(deadline string) (time.Time, error) {
    return time.Parse("2006-01-02 15:04:05", deadline)
}

// Get All Tasks godoc
// @Summary Get all tasks assigned to the user
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {array} models.Task "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/tasks [get]
func GetAllTaskController(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    userID, _ := c.Get("user_id")

    tasks, err := services.GetAllTasksService(db, userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, tasks)
}

// Get Task by ID godoc
// @Summary Get a task by its ID
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Task Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/tasks/{id} [get]
func GetTaskByIDController(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    userID, _ := c.Get("user_id")

    taskID, _ := strconv.Atoi(c.Param("id"))
    task, err := services.GetTaskByIDService(db, uint(taskID), userID.(uint))
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, task)
}

// Get Tasks by Project godoc
// @Summary Get tasks by project ID
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param project_id path int true "Project ID"
// @Success 200 {array} models.Task "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Project Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/projects/{project_id}/tasks [get]
func GetTaskByProjectController(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    userID, _ := c.Get("user_id")

    projectID, _ := strconv.Atoi(c.Param("project_id"))
    tasks, err := services.GetTaskByProjectService(db, uint(projectID), userID.(uint))
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, tasks)
}

// Add Task godoc
// @Summary Add a new task to a project
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param project_id path int true "Project ID"
// @Param input body taskInput true "Task Data"
// @Success 201 {object} models.Task "Task created successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/projects/{project_id}/tasks [post]
func AddTaskController(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    
    projectID, _ := strconv.Atoi(c.Param("project_id"))
    var input taskInput
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    deadline, err := parseDeadline(input.Deadline)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format"})
        return
    }

    task := models.Task{
        Title:       input.Title,
        Description: input.Description,
        Status:      input.Status,
        Deadline:    deadline,
    }

    userID := c.MustGet("user_id").(uint)
    
    if err := services.CreateTaskService(db, uint(projectID), &task, input.AssignedTo, userID); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"data": task})
}

// Update Task godoc
// @Summary Update an existing task
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param project_id path int true "Project ID"
// @Param task_id path int true "Task ID"
// @Param input body taskInput true "Task Data"
// @Success 200 {object} models.Task "Task updated successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Task Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/projects/{project_id}/tasks/{task_id} [put]
func UpdateTaskController(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    userID := c.MustGet("user_id").(uint)
    
    projectID, _ := strconv.Atoi(c.Param("project_id"))
    taskID, _ := strconv.Atoi(c.Param("task_id"))
    
    var input taskInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    deadline, err := parseDeadline(input.Deadline)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format"})
        return
    }

    task := models.Task{
        Title:       input.Title,
        Description: input.Description,
        Status:      input.Status,
        Deadline:    deadline,
    }

    if err := services.UpdateTaskService(db, uint(projectID), uint(taskID), &task, input.AssignedTo, userID); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": task})
}

// Delete Task godoc
// @Summary Delete a task by ID
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string "Task deleted successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Task Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/tasks/{id} [delete]
func DeleteTaskController(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    userID := c.MustGet("user_id").(uint)
    taskID, _ := strconv.Atoi(c.Param("id"))

    if err := services.DeleteTaskService(db, uint(taskID), userID); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}