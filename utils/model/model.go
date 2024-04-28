package model

type PCB struct {
	PID              int
	State            int
	EndState         int
	PC               int
	CPUTime          int
	Quantum           int
	RemainingQuantum int
	Registers        cpuRegister
}

type cpuRegister struct {
	AX int
	BX int
	CX int
	DX int
	EAX int
	EBX int
	ECX int
	EDX int
}

type Instruction struct {
	Operation string
	Param1    string
	Param2    string
}