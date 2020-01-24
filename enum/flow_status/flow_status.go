package flow_status

type FlowStatus int

const (
	Draft = iota //starts with 0
	Inactive
	Active
)
