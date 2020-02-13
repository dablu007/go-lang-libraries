package enum

type FlowType int

const (
	CreditLimit = iota //starts with 0
	Checkout
)

func (flowType FlowType) String() string {
	// declare an array of strings
	flowTypes := [...]string{
		"CreditLimit",
		"Checkout",
	}

	// prevent panicking in case of
	// `flowType` is out of range of Weekday
	if flowType < CreditLimit || flowType > Checkout {
		return "Unknown"
	}
	return flowTypes[flowType]
}
