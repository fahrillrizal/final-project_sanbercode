package models

import "time"

// @model
type Project struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	OwnerID uint `gorm:"not null" json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Collaborators []ProjectCollaborator `gorm:"foreignKey:ProjectID" json:"collaborators"`
}

// @model
type ProjectCollaborator struct {
	ID uint `gorm:"primaryKey" json:"id"`
	ProjectID uint `gorm:"not null;index;constraint:OnDelete:CASCADE" json:"project_id"`
	UserID uint `gorm:"not null" json:"user_id"`
	User User `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt time.Time `json:"created_at"`
}