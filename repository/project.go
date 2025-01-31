package repository

import (
	"PA/models"
	"errors"
	"gorm.io/gorm"
)

func GetAllProjects(db *gorm.DB, userID uint) ([]models.Project, error) {
	var projects []models.Project

	err := db.Preload("Collaborators.User", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, username, email")
    }).
    Where("owner_id = ?", userID).
    Find(&projects).Error

    if err != nil {
        return nil, err
    }
    return projects, nil
}

func GetProjectByID(db *gorm.DB, projectID uint, userID uint) (models.Project, error) {
	var project models.Project

	err := db.Preload("Collaborators.User", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, username, email")
    }).First(&project, projectID).Error

	return project, err
}

func CreateProject(db *gorm.DB, project *models.Project) error {
	return db.Create(project).Error
}

func UpdateProject(db *gorm.DB, project *models.Project, userID uint) error {
    result := db.Model(&models.Project{}).
        Where("id = ? AND owner_id = ?", project.ID, userID).
        Updates(project)
        
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("project not found or unauthorized")
    }
    return nil
}

func DeleteProject(db *gorm.DB, projectID uint, ownerID uint) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Where("project_id = ?", projectID).Delete(&models.ProjectCollaborator{}).Error; err != nil {
            return err
        }
        
        var taskIDs []uint
        if err := tx.Model(&models.Task{}).Where("project_id = ?", projectID).Pluck("id", &taskIDs).Error; err != nil {
            return err
        }
        
        if len(taskIDs) > 0 {
            if err := tx.Where("task_id IN (?)", taskIDs).Delete(&models.TaskAssignment{}).Error; err != nil {
                return err
            }
            if err := tx.Where("project_id = ?", projectID).Delete(&models.Task{}).Error; err != nil {
                return err
            }
        }

        result := tx.Where("id = ? AND owner_id = ?", projectID, ownerID).Delete(&models.Project{})
        if result.Error != nil {
            return result.Error
        }
        if result.RowsAffected == 0 {
            return gorm.ErrRecordNotFound
        }
        return nil
    })
}

func InviteCollaborator(db *gorm.DB, projectID, userID uint) error {
	collab := models.ProjectCollaborator{
		ProjectID: projectID,
		UserID: userID,
	}
	return db.Create(&collab).Error
}

func RemoveCollaborator(db *gorm.DB, projectID, userID uint) error {
    result := db.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&models.ProjectCollaborator{})
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}

func IsOwner(db *gorm.DB, projectID, userID uint) (bool, error) {
    var count int64
    err := db.Model(&models.Project{}).
        Where("id = ? AND owner_id = ?", projectID, userID).
        Count(&count).Error
        
    return count > 0, err
}