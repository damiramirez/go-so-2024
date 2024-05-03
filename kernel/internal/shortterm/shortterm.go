package shortterm

import (
	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/algorithm"
)

func InitShortTermPlani() {

	switch global.KernelConfig.PlanningAlgorithm {
	case "FIFO":
		algorithm.Fifo()
	}

}