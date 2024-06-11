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
	Pid    int `json:"pid"`
	Frames int `json:"frames"`
}

/*type MemAccess struct {
	NumFrame int    `json:"numframe"`
	NumPage  int    `json:"numpage"`
	Offset   int    `json:"offset"`
	Content  uint32 `json:"content"`
	Pid      int    `json:"pid"`
	Largo    int    `json:"largo"`
}*/

type MemStruct struct {
	Pid       int   `json:"pid"`
	Content   int  `json:"content"`
	Length    int   `json:"length"`
	NumFrames []int `json:"numframe"`
	Offset    int   `json:"offset"`
}
