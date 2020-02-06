package response_dto

import (
	uuid "github.com/google/uuid"
)

type JourneyResponseDto struct {
	Name       string                     `json:"name"`
	ExternalId uuid.UUID                  `json:"external_id"`
	Version    string                     `json:"version"`
	Type       string                     `json:"type"`
	Modules    []ModuleVersionResponseDto `json:"modules"`
}

type JourneyResponseDtoList struct {
	Name       string                     `json:"name"`
	ExternalId uuid.UUID                  `json:"external_id"`
	Version    string                     `json:"version"`
	Type       string                     `json:"type"`
	Modules    []ResponseDTO `json:"modules"`
	Sections    []ResponseDTO `json:"sections"`
	Fields    []ResponseDTO `json:"fields"`
}
