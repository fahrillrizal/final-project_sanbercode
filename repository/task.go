package repository

import (
    "PA/models"
    "time"
    "gorm.io/gorm"
)

func GetAllTask(db *gorm.DB, userID uint) ([]models.Task, error) {
    var tasks []models.Task
    err := db.
        Preload("Assignments.User").
        Preload("Project").
        Joins("JOIN projects ON projects.id = tasks.project_id").
        Where("tasks.id IN (SELECT task_id FROM task_assignments WHERE user_id = ?) OR projects.owner_id = ?", userID, userID).
        Find(&tasks).Error
    return tasks, err
}

func GetTaskByID(db *gorm.DB, id uint) (models.Task, error) {
    var task models.Task
    err := db.
        Preload("Assignments.User").
        Preload("Project").
        First(&task, id).Error
    return task, err
}

func GetTaskByProject(db *gorm.DB, projectID uint) ([]models.Task, error) {
    var tasks []models.Task
    err := db.
        Preload("Assignments.User").
        Preload("Project").
        Where("project_id = ?", projectID).
        Find(&tasks).Error
    return tasks, err
}

func CreateTask(db *gorm.DB, task *models.Task, userIDs []uint) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(task).Error; err != nil {
            return err
        }
        for _, userID := range userIDs {
            assignment := models.TaskAssignment{
                TaskID:     task.ID,
                UserID:     userID,
                AssignedAt: time.Now(),
            }
            if err := tx.Create(&assignment).Error; err != nil {
                return err
            }
        }
        return tx.Preload("Assignments", func(db *gorm.DB) *gorm.DB {
            return db.Preload("User")
        }).Preload("Project").First(task).Error
    })
}

func UpdateTask(db *gorm.DB, task *models.Task, userIDs []uint) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(task).Updates(task).Error; err != nil {
            return err
        }
        if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskAssignment{}).Error; err != nil {
            return err
        }
        for _, userID := range userIDs {
            assignment := models.TaskAssignment{
                TaskID:     task.ID,
                UserID:     userID,
                AssignedAt: time.Now(),
            }
            if err := tx.Create(&assignment).Error; err != nil {
                return err
            }
        }
        return tx.Preload("Assignments", func(db *gorm.DB) *gorm.DB {
            return db.Preload("User")
        }).Preload("Project").First(task).Error
    })
}

func DeleteTask(db *gorm.DB, id uint) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Where("task_id = ?", id).Delete(&models.TaskAssignment{}).Error; err != nil {
            return err
        }
        return tx.Delete(&models.Task{}, id).Error
    })
}