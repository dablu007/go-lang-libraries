package response_dto

import "flow/enum"

type FieldVersionsResponseDto struct {
	Name       string
	ExternalId string
	FieldId    int
	IsVisible  bool
	Status     enum.Status
	Version    string
}
