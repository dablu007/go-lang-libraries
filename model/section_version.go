package model

import (
	"time"
)

type SectionVersion struct {
	// gorm.Model
	Id         int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Name       string `gorm:"type:varchar(200)"`
	ExternalId string `gorm:"type:varchar(36)"`
	SectionId  int    `gorm:"column:sectionid"`
	IsVisible  bool
	Version    string
	CreatedOn  time.Time
	DeletedOn  time.Time
}
