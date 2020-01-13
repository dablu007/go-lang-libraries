package model

import "time"

type ModuleVersion struct {
	Id              int
	ModuleId        int
	ExternalId      string
	Version         string
	CreatedOn       time.Time
	DeletedOn       time.Time
	Properties      string
	SectionVersions string
}
