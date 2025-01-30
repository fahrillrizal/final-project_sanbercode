package models

import "time"

type Task struct {
    ID uint `gorm:"primaryKey" json:"id"`
    ProjectID uint `json:"project_id"`
    Project Project `gorm:"foreignKey:ProjectID" json:"project"`
    Title string `json:"title"`
    Description string `json:"description"`
    Status string `json:"status"`
    Assignments []TaskAssignment `gorm:"foreignKey:TaskID" json:"-"`
    AssignedTo []UserResponse `gorm:"-" json:"assigned_to"`
    Deadline time.Time `json:"deadline"`
}

type TaskAssignment struct {
    ID uint `gorm:"primaryKey" json:"id"`
    TaskID uint `json:"task_id"`
    UserID uint `json:"user_id"`
    User User `gorm:"foreignKey:UserID" json:"-"`
    AssignedAt time.Time `json:"assigned_at"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}