package repository

import (
	"gorm.io/gorm"
)

var dbInstance *gorm.DB //TODO: remove global db

func InitRepository(db *gorm.DB) {
	dbInstance = db
}

func GetDB() *gorm.DB {
	return dbInstance
}
