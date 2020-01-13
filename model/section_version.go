package model

import "time"

type SectionVersion struct {
	Id         int
	Name string
	ExternalId string
	SectionId int
	IsVisible bool
	Version string
	CreatedOn  time.Time
	DeletedOn  time.Time
}
