package response_dto

import (
	uuid "github.com/satori/go.uuid"
)

type FieldVersionsResponseDto struct {
	Name       string                 `json:"name"`
	ExternalId uuid.UUID              `json:"external_id"`
	IsVisible  bool                   `json:"is_visible"`
	Version    string                 `json:"version"`
	Properties []PropertryResponseDto `json:"properties"`
}
