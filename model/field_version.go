package model

import "time"

type FieldVersion struct{
	Id int;
	Name string;
	ExternalId string;
	FieldId int;
	IsVisible bool;
	Version string;
	CreatedOn time.Time;
	DeletedOn time.Time;
}