package enum

type ExecutionStrategy int

const (
	Background = iota //starts with 0
	User
)

func (strategy ExecutionStrategy) String() string {
	// declare an array of strings
	strategies := [...]string{
		"Background",
		"User",
	}

	// prevent panicking in case of
	// `flowType` is out of range of Weekday
	if strategy < Background || strategy > User {
		return "Unknown"
	}
	return strategies[strategy]
}
