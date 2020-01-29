package enum

type Status int

const (
	Inactive = iota //starts with 0
	Active
)

func (status Status) String() string {
	// declare an array of strings
	statuses := [...]string{
		"Inactive",
		"Active",
	}

	// prevent panicking in case of
	// `status` is out of range
	if status < Inactive || status > Active {
		return "Unknown"
	}
	return statuses[status]
}
