package response_dto

import (
	"flow/model"
	uuid "github.com/satori/go.uuid"
)

type FieldVersionsResponseDto struct {
	Name       string
	ExternalId uuid.UUID
	IsVisible  bool
	Version    string
	Properties []model.Property
}
