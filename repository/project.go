package repository

import (
	"PA/models"
	"errors"
	"gorm.io/gorm"
)

func GetAllProjects(db *gorm.DB) ([]models.Project, error) {
	var projects []models.Project
	err := db.Find(&projects).Error
	return projects, err
}

func GetProjectByID(db *gorm.DB, projectID uint) (models.Project, error) {
	var project models.Project
	err := db.First(&project, projectID).Error
	return project, err
}

func CreateProject(db *gorm.DB, project *models.Project) error {
	return db.Create(project).Error
}

func UpdateProject(db *gorm.DB, project *models.Project) error {
	return db.Save(project).Error
}

func DeleteProject(db *gorm.DB, projectID uint, ownerID uint) error {
	return db.Where("id = ? AND owner_id = ?", projectID, ownerID).Delete(&models.Project{}).Error
}

func InviteCollaborator(db *gorm.DB, projectID, userID uint) error {
	collab := models.ProjectCollaborator{
		ProjectID: projectID,
		UserID: userID,
	}
	return db.Create(&collab).Error
}

func RemoveCollaborator(db *gorm.DB, projectID, userID uint) error {
	return db.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&models.ProjectCollaborator{}).Error
}

func IsOwner(db *gorm.DB, projectID, userID uint) (bool, error) {
	var project models.Project
	err := db.Where("id = ? AND owner_id = ?", projectID, userID).First(&project).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}