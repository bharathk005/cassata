package repository

import (
	"cassata/models"
)

func GetWorkspaceByID(id string) (*models.Workspace, error) {
	db := GetDB()
	var workspace models.Workspace
	if err := db.First(&workspace, id).Error; err != nil {
		return nil, err
	}
	return &workspace, nil
}

func GetAllWorkspaces() ([]models.Workspace, error) {
	db := GetDB()
	var workspaces []models.Workspace
	if err := db.Find(&workspaces).Error; err != nil {
		return nil, err
	}
	return workspaces, nil
}
