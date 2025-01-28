package services

import (
	"gorm.io/gorm"
	
	"errors"
	"PA/models"
	"PA/repository"
)

func GetAllProjectsService(db *gorm.DB) ([]models.Project, error) {
	return repository.GetAllProjects(db)
}

func GetProjectByIDService(db *gorm.DB, projectID uint) (models.Project, error) {
	var project models.Project
	err := db.First(&project, projectID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Project{}, nil
		}
		return models.Project{}, err
	}

	return project, nil
}

func CreateProjectService(db *gorm.DB, project *models.Project) error {
	return repository.CreateProject(db, project)
}

func UpdateProjectService(db *gorm.DB, project *models.Project, userID uint) error {
	return repository.UpdateProject(db, project)
}

func DeleteProjectService(db *gorm.DB, projectID uint, userID uint) error {
	return repository.DeleteProject(db, projectID, userID)
}

func AddCollaboratorService(db *gorm.DB, projectID, userID, ownerID uint) error {
	isOwner, err := repository.IsOwner(db, projectID, ownerID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("tidak diperbolehkan karena anda bukan owner")
	}

	return repository.InviteCollaborator(db, projectID, userID)
}

func RemoveCollaboratorService(db *gorm.DB, projectID, userID, ownerID uint) error {
	isOwner, err := repository.IsOwner(db, projectID, ownerID)

	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("tidak diperbolehkan karena anda bukan owner")
	}

	return repository.RemoveCollaborator(db, projectID, userID)
}