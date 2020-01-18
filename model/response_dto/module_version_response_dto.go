package response_dto

import (
	"flow/model"
	uuid "github.com/satori/go.uuid"
)

type ModuleVersionResponseDto struct {
	Name       string
	Version    string
	ExternalId uuid.UUID
	Properties []model.Property
	Sections   []SectionVersionsResponseDto
}
