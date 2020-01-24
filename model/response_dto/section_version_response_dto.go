package response_dto

import (
	uuid "github.com/satori/go.uuid"
)

type SectionVersionsResponseDto struct {
	Name       string                     `json:"name"`
	ExternalId uuid.UUID                  `json:"external_id"`
	IsVisible  bool                       `json:"is_visible"`
	Version    string                     `json:"verison"`
	Properties []PropertryResponseDto     `json:"properties"`
	Fields     []FieldVersionsResponseDto `json:"fields"`
}
