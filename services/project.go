package services

import (
	"gorm.io/gorm"
	
	"errors"
	"PA/models"
	"PA/repository"
)

func GetAllProjectsService(db *gorm.DB, userID uint) ([]models.Project, error) {
    var projects []models.Project

    err := db.Where("owner_id = ?", userID).Find(&projects).Error
    if err != nil {
        return nil, err
    }

    var collaboratorProjects []models.Project
    err = db.Joins("JOIN project_collaborators ON project_collaborators.project_id = projects.id").
        Where("project_collaborators.user_id = ?", userID).
        Find(&collaboratorProjects).Error
    if err != nil {
        return nil, err
    }

    projects = append(projects, collaboratorProjects...)

    return projects, nil
}

func GetProjectByIDService(db *gorm.DB, projectID uint, userID uint) (models.Project, error) {
    var project models.Project

    err := db.Where("id = ? AND (owner_id = ? OR id IN (SELECT project_id FROM project_collaborators WHERE user_id = ?))", projectID, userID, userID).
        First(&project).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return models.Project{}, errors.New("anda tidak memiliki akses untuk project ini")
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