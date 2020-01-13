package model

import "time"
import "flow/enum"

type Module struct {
	Id        int
	Name      string
	Status    enum.Status
	IsVisible bool
	TenantId  string
	CreatedOn time.Time
	DeletedOn time.Time
}
