package model

import (
	"flow/enum"
	"time"
)

type Field struct {
	// gorm.Model
	Id            int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Name          string `gorm:"type:varchar(200)"`
	SectionId     int    `gorm:"column:sectionid"`
	Status        enum.Status
	IsVisible     bool
	TenantId      string
	CreatedOn     time.Time
	DeletedOn     time.Time
	Section       Section `gorm:"foreignkey:fk_fields_sectionid"`
	FieldVersions []FieldVersion
}
