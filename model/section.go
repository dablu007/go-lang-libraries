package model

import (
	"flow/enum"
	"time"
)

type Section struct {
	// gorm.Model
	Id              int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Name            string `gorm:"type:varchar(200)"`
	Status          enum.Status
	IsVisible       bool
	TenantId        string
	CreatedOn       time.Time
	DeletedOn       time.Time
	Module          Module `gorm:"foreignkey:fk_sections_moduleid"`
	Fields          []Field
	SectionVersions []SectionVersion
}
