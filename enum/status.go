package enum

type Status int

const (
	Pending Status = iota
	Invalid
	Partial
	Expired
	Completed
)

func (s Status) String() string {
	return [...]string{"Pending", "Invalid", "Partial", "Expired", "Completed"}[s]
}
