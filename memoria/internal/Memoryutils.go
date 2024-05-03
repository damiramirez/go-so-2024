package internal

type ProcessPath struct {
	Path string `json:"path"`
	Pid  int    `json:"pid"`
}
type ProcessAssets struct {
	Pc  int `json:"pc"`
	Pid int `json:"pid"`
}
type ProcessDelete struct {
	Pid int `json:"pid"`
}
type Page struct {
	PageNumber int `json:"page_number"`
	Pid        int `json:"pid"`
}
type Resize struct {
	Tipo   string `json:"type"`
	Pid    int    `json:"pid"`
	Frames int    `json:"frames"`
}

type MemAccess struct {
	Tipo          string `json:"type"`
	Adress        int    `json:"adress"`
	NumberOfPages int    `json:"numberofpages"`
}

