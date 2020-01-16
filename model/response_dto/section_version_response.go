package response_dto

import "flow/enum"

type SectionVersionsResponseDto struct {
	Name       string
	ExternalId string
	IsVisible  bool
	Version    string
	Status     enum.Status
	Fields     []FieldVersionsResponseDto
}
