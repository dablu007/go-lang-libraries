package response_dto

import (
	uuid "github.com/google/uuid"
)

type FieldVersionsResponseDto struct {
	Name       string                 `json:"name"`
	ExternalId uuid.UUID              `json:"externalId"`
	IsVisible  bool                   `json:"isVisible"`
	Version    string                 `json:"version"`
	Properties []PropertryResponseDto `json:"properties"`
}

type FieldVersionDTO struct {
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	ExternalId uuid.UUID `json:"externalId"`
}
