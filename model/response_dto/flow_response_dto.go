package response_dto

import (
	uuid "github.com/satori/go.uuid"
)

type FlowResponseDto struct {
	Name       string                     `json:"name"`
	ExternalId uuid.UUID                  `json:"external_id"`
	Version    string                     `json:"version"`
	Type       string                     `json:"type"`
	Modules    []ModuleVersionResponseDto `json:"modules"`
}
