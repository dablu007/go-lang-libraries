package model

import "time"

type Section struct {
	Id        int
	Name      string
	ModuleId  int
	IsVisible bool
	Version   string
	CreatedOn time.Time
	DeletedOn time.Time
}
