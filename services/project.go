package services

import (
	"gorm.io/gorm"
	
	"errors"
	"PA/models"
	"PA/repository"
)

func GetAllProjectsService(db *gorm.DB, userID uint) ([]models.Project, error) {
	return repository.GetAllProjects(db, userID)
}

func GetProjectByIDService(db *gorm.DB, projectID uint, userID uint) (models.Project, error) {
	project, err := repository.GetProjectByID(db, projectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Project{}, errors.New("project tidak ditemukan")
		}
		return models.Project{}, err
	}

	isCollaborator := false
	for _, collaborator := range project.Collaborators {
		if collaborator.UserID == userID {
			isCollaborator = true
			break
		}
	}

	if project.OwnerID != userID && !isCollaborator {
		return models.Project{}, errors.New("anda tidak memiliki akses ke project ini")
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
        if userID != ownerID {
            return errors.New("tidak diperbolehkan karena anda bukan owner atau tidak dapat menghapus orang lain selain diri anda sendiri")
        }
    }

    return repository.RemoveCollaborator(db, projectID, userID)
}