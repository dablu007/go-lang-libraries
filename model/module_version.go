package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ModuleVersion struct {
	gorm.Model
	Id              int `gorm:"primary_key";"AUTO_INCREMENT"`
	ModuleId        int
	ExternalId      string `gorm:"type:varchar(36)"`
	Version         string
	CreatedOn       time.Time
	DeletedOn       time.Time
	Properties      string
	SectionVersions string
}
