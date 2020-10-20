package masa

import (
	"unsafe"
)

const (
	MemoryAllocated      = "allocated"
	MemoryAllocating     = "allocating"
	MemoryAllocateFailed = "failed"
	MemoryWaitToAllocate = "waiting"

	Tiny   = 16
	Little = 32

	Tiny_   = "tiny"
	Little_ = "little"
	Enough_ = "enough"
)

type emptymem struct {
	ptr unsafe.Pointer
}

// memory struct
type masa struct {
	ptr     unsafe.Pointer
	offset  uint
	size    uintptr
	status  string
	memList map[int]*emptymem
}

type Masa interface {
	init() *masa
	getAllocationSize() uintptr
	level() string
	allocate(size uintptr)
	isOutOfMemory(size uintptr) bool
	doAllocLittle(size uintptr) (status string)
	doAllocTiny(size uintptr) (status string)
	doAllocEnough(size uintptr) (status string)
}

func (m *masa) init() *masa {
	m = &masa{}
	return m
}

func (m *masa) getAllocationSize() uintptr {
	return m.size
}

func (m *masa) level(size uintptr) string {
	if size <= Little {
		if size >= Tiny {
			return Little_
		}
		return Tiny_
	}
	return Enough_
}

func (m *masa) allocate(size uintptr) {
	switch m.level(size) {
	case Little_:
		m.doAllocLittle(size)
	case Tiny_:
		m.doAllocTiny(size)
	case Enough_:
		m.doAllocEnough(size)
	}
}

func (m *masa) isOutOfMemory(size uintptr) bool {
	if int(size) < Tiny {
		cnt := 0
		for i := 0; i < int(size); i++ {
			if m.memList[i] != nil {
				cnt++
			}
		}
		if cnt < int(size) {
			return true
		}
		return false
	} else if int(size) >= Tiny && int(size) <= Little {
		cnt := 0
		for i := Tiny; i <= Little; i++ {
			if m.memList[i] != nil {
				cnt++
			}
		}
		if cnt < int(size) {
			return true
		}
		return false
	} else {
		var i int
		cnt := 0
		for ; i > Little; i++ {
			if m.memList[i] != nil {
				cnt++
			}
		}
		if cnt < int(size) {
			return true
		}
		return false
	}
}

func (m *masa) doAllocLittle(size uintptr) (status string) {
	if m.isOutOfMemory(size) {
		return MemoryAllocateFailed
	}
	// TODO: allocation
	return MemoryAllocated
}

func (m *masa) doAllocTiny(size uintptr) (status string) {
	if m.isOutOfMemory(size) {
		return MemoryAllocateFailed
	}
	// TODO: allocation
	return MemoryAllocated
}

func (m *masa) doAllocEnough(size uintptr) (status string) {
	if m.isOutOfMemory(size) {
		return MemoryAllocateFailed
	}
	// TODO: allocation
	return MemoryAllocated
}
