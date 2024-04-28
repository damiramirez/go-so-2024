package internal

import (
	"testing"

	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func TestSet(t *testing.T) {
	pcb := &model.PCB{
		Registers: model.CPURegister{},
	}
	instruction := model.Instruction{
		Operation:  "SET",
		Parameters: []string{"AX", "100"},
	}

	set(pcb, instruction)

	expected := 100
	if pcb.Registers.AX != expected {
		t.Errorf("Set failed: expected AX = %d, got %d", expected, pcb.Registers.AX)
	}
}

func TestSum(t *testing.T) {
	pcb := &model.PCB{
		Registers: model.CPURegister{
			AX: 10,
			BX: 5,
		},
	}
	instruction := model.Instruction{
		Operation:  "SUM",
		Parameters: []string{"AX", "BX"},
	}

	sum(pcb, instruction)

	expected := 15
	if pcb.Registers.AX != expected {
		t.Errorf("Sum failed: expected AX = %d, got %d", expected, pcb.Registers.AX)
	}
}

func TestSub(t *testing.T) {
	pcb := &model.PCB{
		Registers: model.CPURegister{
			AX: 10,
			BX: 5,
		},
	}
	instruction := model.Instruction{
		Operation:  "SUB",
		Parameters: []string{"AX", "BX"},
	}

	sub(pcb, instruction)

	expected := 5
	if pcb.Registers.AX != expected {
		t.Errorf("Sub failed: expected AX = %d, got %d", expected, pcb.Registers.AX)
	}
}
