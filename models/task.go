package models

import "time"

type Task struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ProjectID   uint   `json:"project_id"`
	Project     Project `gorm:"foreignKey:ProjectID" json:"project"` // Memisahkan gorm dan json dengan spasi
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	AssignedToID uint   `json:"assigned_to_id"`
	AssignedTo  User   `gorm:"foreignKey:AssignedToID" json:"assigned_to"` // Memisahkan gorm dan json dengan spasi
	Deadline    time.Time `json:"deadline"`
}

type TaskAssignment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `json:"task_id"`
	UserID     uint      `json:"user_id"`
	AssignedAt time.Time `json:"assigned_at"`
	Role       string    `json:"role"`
}