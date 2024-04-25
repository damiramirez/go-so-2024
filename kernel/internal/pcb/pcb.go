package pcb

import "github.com/sisoputnfrba/tp-golang/kernel/global"

const (
	NEW = 1
	READY = 2
	EXEC = 3
	BLOCKED = 4
	EXIT = 5
)

type PCB struct {
	PID              int
	State            int
	EndState         int
	PC               int
	CPUTime          int
	Quantum           int
	RemainingQuantum int
	Registers        CpuRegisters
}

type CpuRegisters struct {
	AX int
	BX int
	CX int
	DX int
	EAX int
	EBX int
	ECX int
	EDX int
}

func CreateNewProcess() *PCB {
	return &PCB{
		PID:     global.GetNextPID(),
		State:   NEW,
		Quantum:  global.KernelConfig.Quantum,
		RemainingQuantum: global.KernelConfig.Quantum,
	}
}
