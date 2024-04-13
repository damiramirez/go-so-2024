package global

type Config struct {
	Port             int    `json:"port"`
	IPKernel         string `json:"ip_kernel"`
	PortKernel       int    `json:"port_kernel"`
	IPMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	NumberFellingTLB int    `json:"number_felling_tlb"`
	AlgorithmTLB     string `json:"algorithm_tlb"`
}

var CPUConfig *Config
