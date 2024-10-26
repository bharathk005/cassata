package models

import (
	"gorm.io/gorm"

	"cassata/utils"
)

type User struct {
	gorm.Model
	Name         string
	PasswordHash string
	Workspaces   []Workspace `gorm:"many2many:user_workspace;"`
}

type Workspace struct {
	gorm.Model
	Name        string
	Users       []User       `gorm:"many2many:user_workspace;"`
	Permissions []Permission `gorm:"many2many:workspace_permissions;"`
}

type Permission struct {
	gorm.Model
	Resource string
	Action   utils.Verb
}

type ResourceMap struct {
	gorm.Model
	Provider      string
	ResourceGroup string
	ResourceType  string
	K8sApiGroup   string
	K8sApiVersion string
	K8sResource   string
}
