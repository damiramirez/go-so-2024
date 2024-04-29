package internal

type ProcessPath struct {
	Path string `json:"path"`
	Pid  int    `json:"pid"`
}
type PCB struct {
	Pc  int `json:"pc"`
	Pid int `json:"pid"`
}
type MemoryST struct {
	spaces []byte
}
type PageTable struct {
	pages []byte
}

