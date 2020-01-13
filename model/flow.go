package model

import (
	"flow/enum"
	"time"
)

type Flow struct {
	Id             int
	Name           string
	Version        string
	Type           enum.FlowType
	Status         enum.Status
	MerchantId     string
	ModuleVersions string
	CreatedOn      time.Time
	DeletedOn      time.Time
}
