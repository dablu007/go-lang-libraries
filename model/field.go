package model

import "time"

type Field struct {
	Id        int
	Name      string
	SectionId int
	IsVisible bool
	Version   string
	CreatedOn time.Time
	DeletedOn time.Time
}
