package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Section struct {
	Id        int       `gorm:"primary_key";"AUTO_INCREMENT";"column:id"`
	Name      string    `gorm:"type:varchar(200)";"column:name"`
	TenantId  uuid.UUID `gorm:"column:tenant_id"`
	CreatedOn time.Time `gorm:"column:created_on"`
	DeletedOn time.Time `gorm:"column:deleted_on"`
}
