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

type taskInput struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Status      string   `json:"status"`
    AssignedTo  []uint   `json:"assigned_to"`
    Deadline    string   `json:"deadline"`
}

func parseDeadline(deadline string) (time.Time, error) {
    return time.Parse("2006-01-02 15:04:05", deadline)
}

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