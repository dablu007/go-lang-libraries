package response_dto

import "flow/enum"

type ModuleVersionResponseDto struct {
	Name       string
	Version    string
	ExternalId string
	Status     enum.Status
	Properties string
	Sections   []SectionVersionsResponseDto
}
