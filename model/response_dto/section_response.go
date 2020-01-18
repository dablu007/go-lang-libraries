package response_dto

import "flow/enum"

type SectionResponseDto struct {
	Name       string
	ExternalId string
	IsVisible  bool
	Version    string
	Status     enum.Status
	Fields     []FieldVersionsResponseDto
}
