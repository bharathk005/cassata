package repository

import (
	"cassata/models"
)

func GetUserByID(id string) (*models.User, error) {
	db := GetDB()
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
