package model

import (
	"flow/enum"
	"flow/enum/flow_status"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Flow struct {
	// gorm.Model
	Id             int                    `gorm:"primary_key";"AUTO_INCREMENT";"column:id"`
	ExternalId     uuid.UUID              `gorm:"column:external_id"`
	Name           string                 `gorm:"column:name"`
	Version        string                 `gorm:"column:version"`
	Type           enum.FlowType          `gorm:"column:type"`
	Status         flow_status.FlowStatus `gorm:"column:status"`
	FlowContext    string                 `gorm:"column:flow_context";"type:json"`
	ModuleVersions string                 `gorm:"column:module_versions";"type:json"`
	CreatedOn      time.Time              `gorm:"column:created_on"`
	DeletedOn      time.Time              `gorm:"column:deleted_on"`
}
