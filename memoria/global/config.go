package global

type MemoryConfig struct {
	Port             int    `json:"port"`
	MemorySize       int    `json:"memory_size"`
	PageSize         int    `json:"page_size"`
	InstructionsPath string `json:"instructions_path"`
	DelayResponse    int    `json:"delay_response"`
}
