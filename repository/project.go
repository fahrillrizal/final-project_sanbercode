package repository

import (
	"PA/models"
	"errors"
	"gorm.io/gorm"
)

func GetAllProjects(db *gorm.DB, userID uint) ([]models.Project, error) {
	var projects []models.Project

	err := db.Preload("Collaborators").Where("owner_id = ?", userID).Find(&projects).Error
	if err != nil {
		return nil, err
	}

	var collaboratorProjects []models.Project
	err = db.Joins("JOIN project_collaborators ON project_collaborators.project_id = projects.id").
		Preload("Collaborators").
		Where("project_collaborators.user_id = ?", userID).
		Find(&collaboratorProjects).Error
	if err != nil {
		return nil, err
	}

	projects = append(projects, collaboratorProjects...)
	return projects, nil
}

func GetProjectByID(db *gorm.DB, projectID uint, userID uint) (models.Project, error) {
	var project models.Project

	err := db.Preload("Collaborators").
		Where("id = ? AND (owner_id = ? OR id IN (SELECT project_id FROM project_collaborators WHERE user_id = ?))", projectID, userID, userID).
		First(&project).Error

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