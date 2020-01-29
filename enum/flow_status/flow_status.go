package flow_status

type FlowStatus int

const (
	Draft   = iota //starts from 0
	Inactive
	Active     
)

func (flowStatus FlowStatus) String() string {
	// declare an array of strings
	flowStatuses := [...]string{
		"Draft",
		"Inactive",
		"Active",
	}

	// prevent panicking in case of
	// `flowStatus` is out of range
	if flowStatus < Draft || flowStatus > Active {
		return "Unknown"
	}
	return flowStatuses[flowStatus]
}
