package model

import (
	"flow/enum"
	"time"
)

type Module struct {
	// gorm.Model
	Id             int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Name           string `gorm:"type:varchar(200)"`
	Status         enum.Status
	IsVisible      bool
	TenantId       string
	CreatedOn      time.Time
	DeletedOn      time.Time
	ModuleVersions []ModuleVersion
}
