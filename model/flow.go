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
	Name           string                 `gorm:"column:name"`
	Version        string                 `gorm:"column:version"`
	Type           enum.FlowType          `gorm:"column:type"`
	Status         flow_status.FlowStatus `gorm:"column:status"`
	MerchantId     uuid.UUID              `gorm:"column:merchantid"`
	ModuleVersions string                 `gorm:"column:moduleversions"`
	CreatedOn      time.Time              `gorm:"column:createdon"`
	DeletedOn      time.Time              `gorm:"column:deletedon"`
}
