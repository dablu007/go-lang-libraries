package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type FieldVersion struct {
	gorm.Model
	Id         int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Name       string `gorm:"type:varchar(200)"`
	ExternalId string `gorm:"type:varchar(36)"`
	FieldId    int
	IsVisible  bool
	Version    string
	CreatedOn  time.Time
	DeletedOn  time.Time
}
