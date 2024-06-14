package tlb

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

type TLBEntry struct {
	PID    int
	Page   int
	Frame  int
	Access int // LRU
}

type TLB struct {
	entries      []TLBEntry
	capacity     int
	replacement  string
	accessCounter int // LRU
	fifoPointer     int // Para FIFO
}

func NewTLB(capacity int, replacement string) *TLB {
	return &TLB{
			entries:     make([]TLBEntry, 0, capacity),
			capacity:    capacity,
			replacement: replacement,
	}
}

func (tlb *TLB) Search(pid, page int) (int, bool) {
	for i, entry := range tlb.entries {
		if entry.PID == pid && entry.Page == page {
			if tlb.replacement == "LRU" {
				tlb.entries[i].Access = tlb.accessCounter
				tlb.accessCounter++
			}
			global.Logger.Log(fmt.Sprintf("Encontre %d pagina - Frame %d", page, entry.Frame), log.DEBUG)
			return entry.Frame, true // TLB Hit
		}
	}
	global.Logger.Log(fmt.Sprintf("No encontre %d pagina", page), log.DEBUG)

	return -1, false // TLB Miss
}


func (tlb *TLB) AddEntry(pid, page, frame int) {
	if len(tlb.entries) >= tlb.capacity {
		tlb.replaceEntry(pid, page, frame)
		global.Logger.Log(fmt.Sprintf("Se remplazo pagina %d -> frame %d", page, frame), log.DEBUG)
	} else {
		tlb.entries = append(tlb.entries, TLBEntry{
			PID:    pid,
			Page:   page,
			Frame:  frame,
			Access: tlb.accessCounter,
		})
		tlb.accessCounter++
		global.Logger.Log(fmt.Sprintf("Se agrego pagina %d -> frame %d", page, frame), log.DEBUG)
	}
}


func (tlb *TLB) replaceEntry(pid, page, frame int) {
	var index int
	if tlb.replacement == "FIFO" {
		index = tlb.fifoPointer
		tlb.fifoPointer = (tlb.fifoPointer + 1) % tlb.capacity
	} else if tlb.replacement == "LRU" {
		index = tlb.findLRUIndex()
	}

	tlb.entries[index] = TLBEntry{
		PID:    pid,
		Page:   page,
		Frame:  frame,
		Access: tlb.accessCounter,
	}
	tlb.accessCounter++
}


func (tlb *TLB) findLRUIndex() int {
	lruIndex := 0
	minAccess := tlb.entries[0].Access
	for i, entry := range tlb.entries {
		if entry.Access < minAccess {
			lruIndex = i
			minAccess = entry.Access
		}
	}
	return lruIndex
}

