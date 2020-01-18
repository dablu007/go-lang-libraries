package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type ModuleVersion struct {
	Id              int       `gorm:"primary_key";"AUTO_INCREMENT";"column:id"`
	Name            string    `gorm:"column:name"`
	ModuleId        int       `gorm:"column:module_id"`
	ExternalId      uuid.UUID `gorm:"column:external_id"`
	Version         string    `gorm:"column:version"`
	CreatedOn       time.Time `gorm:"column:created_on"`
	DeletedOn       time.Time `gorm:"column:deleted_on"`
	Properties      string    `gorm:"column:properties";"type:json"`
	SectionVersions string    `gorm:"column:section_versions";"type:json"`
}
