package response_dto

import (
	"flow/enum"

	uuid "github.com/satori/go.uuid"
)

type FlowResponseDto struct {
	Name       string                     `json:"name"`
	ExternalId uuid.UUID                  `json:"external_id"`
	Version    string                     `json:"version"`
	Type       enum.FlowType              `json:"type"`
	Modules    []ModuleVersionResponseDto `json:"modules"`
}
