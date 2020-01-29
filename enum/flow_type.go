package enum

type FlowType int

const (
	CL = iota //starts with 0
	PL
)

func (flowType FlowType) String() string {
	// declare an array of strings
	flowTypes := [...]string{
		"CL",
		"PL",
	}

	// prevent panicking in case of
	// `flowType` is out of range of Weekday
	if flowType < CL || flowType > PL {
		return "Unknown"
	}
	return flowTypes[flowType]
}
