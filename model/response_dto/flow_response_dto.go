package response_dto

import (
	"flow/enum"

	uuid "github.com/satori/go.uuid"
)

type FlowResponseDto struct {
	Name       string
	ExternalId uuid.UUID
	Version    string
	Type       enum.FlowType
	Modules    []ModuleVersionResponseDto
}
