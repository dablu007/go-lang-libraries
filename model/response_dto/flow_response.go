package response_dto

import "flow/enum"
import "flow/enum/flow_status"

type FlowResponseDto struct {
	Name    string
	Version string
	Type    enum.FlowType
	Status  flow_status.FlowStatus
	Modules []ModuleVersionResponseDto
}
