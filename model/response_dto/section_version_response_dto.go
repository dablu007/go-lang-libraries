package response_dto

import (
	uuid "github.com/google/uuid"
)

type SectionVersionsResponseDto struct {
	Name       string                     `json:"name"`
	ExternalId uuid.UUID                  `json:"externalId"`
	IsVisible  bool                       `json:"isVisible"`
	Version    string                     `json:"verison"`
	Properties []PropertryResponseDto     `json:"properties"`
	Fields     []FieldVersionsResponseDto `json:"fields"`
}
