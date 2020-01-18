package model

import (
	"time"
)

type ModuleVersion struct {
	Id         int    `gorm:"primary_key";"AUTO_INCREMENT"`
	ModuleId   int    `gorm:"column:moduleid"`
	ExternalId string `gorm:"type:varchar(36)"`
	Version    string
	CreatedOn  time.Time
	DeletedOn  time.Time
	Properties string
	Sections   string `gorm:"column:sections"`
}
