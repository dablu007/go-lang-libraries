package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type SectionVersion struct {
	gorm.Model
	Id         int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Name       string `gorm:"type:varchar(200)"`
	ExternalId string `gorm:"type:varchar(36)"`
	SectionId  int
	IsVisible  bool
	Version    string
	CreatedOn  time.Time
	DeletedOn  time.Time
}
