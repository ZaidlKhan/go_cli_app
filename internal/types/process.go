package types

type Process struct {
	PID         int
	ProcessName string
	Address     string
	Port        uint16
	URL         string
	GitBranch   string
}
