package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type FieldVersion struct {
	Id         int       `gorm:"primary_key";"AUTO_INCREMENT";"column:id"`
	Name       string    `gorm:"type:varchar(200)";"column:name"`
	ExternalId uuid.UUID `gorm:"column:external_id"`
	FieldId    int       `gorm:"column:field_id"`
	IsVisible  bool      `gorm:"column:is_visible"`
	Version    string    `gorm:"column:version"`
	Properties string    `gorm:"column:properties";"type:json"`
	CreatedOn  time.Time `gorm:"column:created_on"`
	DeletedOn  time.Time `gorm:"column:deleted_on"`
}
