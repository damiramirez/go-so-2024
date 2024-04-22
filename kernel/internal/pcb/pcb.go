package pcb

import "github.com/sisoputnfrba/tp-golang/kernel/global"

type PCB struct {
	PID              int
	State            string
	EndState         string
	PC               int
	CPUTime          int
	Quatum           int
	RemainingQuantum int
	Registers        *CpuRegisters
}

type CpuRegisters struct {
	AX int
	BX int
	CX int
	DX int
}

func CreateNewProcess() *PCB {
	return &PCB{
		PID:     global.GetNextPID(),
		State:   "NEW",
		PC:      0,
		CPUTime: 0,
		Quatum:  global.KernelConfig.Quantum,
		Registers: &CpuRegisters{
			AX: 0,
			BX: 0,
			CX: 0,
			DX: 0,
		},
		RemainingQuantum: global.KernelConfig.Quantum,
		EndState:         "",
	}
}
