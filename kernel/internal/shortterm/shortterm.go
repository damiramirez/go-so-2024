package shortterm

import (
	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/algorithm"
)

func initShortTermPlani() {

	switch global.KernelConfig.PlanningAlgorithm {
	case "FIFO":
		go algorithm.Fifo()
	}

}