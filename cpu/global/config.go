package global

type Config struct {
	Port             int    `json:"port"`
	IPMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	NumberFellingTLB int    `json:"number_felling_tlb"`
	AlgorithmTLB     string `json:"algorithm_tlb"`
}
