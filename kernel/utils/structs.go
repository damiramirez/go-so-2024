package utils

type ProcessPID struct {
	PID int `json:"pid"`
}

type ProcessState struct {
	PID   int    `json:"pid"`
	State string `json:"state"`
}

type ProcessPath struct {
	Path string `json:"path"`
}
type NewDevice struct {
	Port  int    `json:"port"`
	Usage bool   `json:"usage"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}
