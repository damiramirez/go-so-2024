package model

type PCB struct {
	PID              int
	State            string
	PC               int
	CPUTime          int
	Quantum          int
	RemainingQuantum int
	DisplaceReason   string
	Registers        CPURegister
	Instruction      Instruction
}

type CPURegister struct {
	AX  int
	BX  int
	CX  int
	DX  int
	EAX int
	EBX int
	ECX int
	EDX int
	SI  int
	DI  int
}

type Instruction struct {
	Operation  string
	Parameters []string
}

type ProcessInstruction struct {
	Pc  int `json:"pc"`
	Pid int `json:"pid"`
}
