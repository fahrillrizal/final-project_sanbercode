package services

import (
    "errors"
    "PA/models"
    "PA/repository"
    "gorm.io/gorm"
)

func mapAssignments(task *models.Task) {
    task.AssignedTo = make([]models.UserResponse, len(task.Assignments))
    for i, assignment := range task.Assignments {
        task.AssignedTo[i] = models.UserResponse{
            ID:       assignment.User.ID,
            Username: assignment.User.Username,
            Email:    assignment.User.Email,
        }
    }
}

func GetAllTasksService(db *gorm.DB, userID uint) ([]models.Task, error) {
    tasks, err := repository.GetAllTask(db, userID)
    if err != nil {
        return nil, err
    }
    for i := range tasks {
        mapAssignments(&tasks[i])
    }
    return tasks, nil
}

func GetTaskByIDService(db *gorm.DB, id, userID uint) (models.Task, error) {
    task, err := repository.GetTaskByID(db, id)
    if err != nil {
        return models.Task{}, err
    }

    if task.Project.OwnerID != userID {
        isAssigned := false
        for _, assignment := range task.Assignments {
            if assignment.UserID == userID {
                isAssigned = true
                break
            }
        }
        if !isAssigned {
            return models.Task{}, errors.New("unauthorized access")
        }
    }
    
    mapAssignments(&task)
    return task, nil
}

func GetTaskByProjectService(db *gorm.DB, projectID, userID uint) ([]models.Task, error) {
    project, err := repository.GetProjectByID(db, projectID, userID)
    if err != nil {
        return nil, err
    }

    isAuthorized := project.OwnerID == userID
    if !isAuthorized {
        for _, collab := range project.Collaborators {
            if collab.UserID == userID {
                isAuthorized = true
                break
            }
        }
    }
    if !isAuthorized {
        return nil, errors.New("unauthorized access")
    }

    tasks, err := repository.GetTaskByProject(db, projectID)
    for i := range tasks {
        mapAssignments(&tasks[i])
    }
    return tasks, err
}

func validateUsersInProject(db *gorm.DB, projectID uint, userIDs []uint) error {
    project, err := repository.GetProjectByID(db, projectID, 0)
    if err != nil {
        return err
    }

    validUsers := map[uint]bool{project.OwnerID: true}
    for _, collab := range project.Collaborators {
        validUsers[collab.UserID] = true
    }

    for _, uid := range userIDs {
        if !validUsers[uid] {
            return errors.New("invalid user assignment")
        }
    }
    return nil
}

func CreateTaskService(db *gorm.DB, projectID uint, task *models.Task, userIDs []uint, currentUserID uint) error {
    var project models.Project
    if err := db.First(&project, projectID).Error; err != nil {
        return errors.New("project tidak ditemukan")
    }

    isOwner, err := repository.IsOwner(db, projectID, currentUserID)
    if err != nil {
        return err
    }
    
    isCollaborator := false
    if !isOwner {
        project, _ := repository.GetProjectByID(db, projectID, 0)
        for _, collab := range project.Collaborators {
            if collab.UserID == currentUserID {
                isCollaborator = true
                break
            }
        }
    }
    
    if !isOwner && !isCollaborator {
        return errors.New("hanya owner/collaborator yang bisa membuat task")
    }

    if err := validateUsersInProject(db, projectID, userIDs); err != nil {
        return err
    }
    
    task.ProjectID = projectID
    if err := repository.CreateTask(db, task, userIDs); err != nil {
        return err
    }
    mapAssignments(task)
    return nil
}

func UpdateTaskService(db *gorm.DB, projectID, taskID uint, task *models.Task, userIDs []uint, userID uint) error {
    isOwner, err := repository.IsOwner(db, projectID, userID)
    if err != nil {
        return err
    }
    
    isCollaborator := false
    if !isOwner {
        project, _ := repository.GetProjectByID(db, projectID, 0)
        for _, collab := range project.Collaborators {
            if collab.UserID == userID {
                isCollaborator = true
                break
            }
        }
    }
    
    if !isOwner && !isCollaborator {
        return errors.New("unauthorized access")
    }

    if err := validateUsersInProject(db, projectID, userIDs); err != nil {
        return err
    }
    
    task.ID = taskID
    task.ProjectID = projectID
    
    if err := repository.UpdateTask(db, task, userIDs); err != nil {
        return err
    }
    
    mapAssignments(task)
    return nil
}

func DeleteTaskService(db *gorm.DB, id uint, userID uint) error {
    task, err := repository.GetTaskByID(db, id)
    if err != nil {
        return err
    }
    
    isOwner, err := repository.IsOwner(db, task.ProjectID, userID)
    if err != nil {
        return err
    }
    
    if !isOwner {
        return errors.New("unauthorized: hanya owner yang bisa menghapus task")
    }
    
    return repository.DeleteTask(db, id)
}