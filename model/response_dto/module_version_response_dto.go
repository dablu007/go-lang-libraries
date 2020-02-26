package response_dto

import (
	uuid "github.com/google/uuid"
)

type ModuleVersionResponseDto struct {
	Name              string                       `json:"name"`
	Version           string                       `json:"version"`
	ExternalId        uuid.UUID                    `json:"externalId"`
	Properties        []PropertryResponseDto       `json:"properties"`
	ExecutionStrategy string                       `json:"executionStrategy"`
	Sections          []SectionVersionsResponseDto `json:"sections"`
}

type ResponseDTO struct {
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	ExternalId uuid.UUID `json:"externalId"`
}

type ModuleVersionResponseWithExecutionStrategyDTO struct {
	Name              string    `json:"name"`
	Version           string    `json:"version"`
	ExternalId        uuid.UUID `json:"externalId"`
	ExecutionStrategy string    `json:"executionStrategy"`
}
