package masa

import (
	"unsafe"

	log "github.com/sirupsen/logrus"
)

const (
	MemoryInitializing   = "init"
	MemoryAllocated      = "allocated"
	MemoryAllocating     = "allocating"
	MemoryAllocateFailed = "failed"
	MemoryWaitToAllocate = "waiting"

	FullOfWaitList = "full"

	Tiny   = 16
	Little = 32

	Tiny_   = "tiny"
	Little_ = "little"
	Enough_ = "enough"
)

type emptymem struct {
	ptr unsafe.Pointer
}

type List struct {
	ptr      unsafe.Pointer
	size     int
	priority int
}

type Task struct {
}

// memory struct
type masa struct {
	ptr     unsafe.Pointer
	offset  uint
	size    uintptr
	status  string
	memList map[int]*emptymem
	cached  *List
	task    *Task
	v       interface{}
}

type Masa interface {
	init() *masa
	getAllocationSize() uintptr
	level() string
	allocate(size uintptr, task ...*Task)
	sliceChecking(s interface{}) (bool, interface{})
	checkAllocationStatus(size uintptr, task ...*Task) string
	isOutOfMemory(size uintptr) bool
	check(size uintptr, task ...*Task)
	doAllocLittle(size uintptr, task ...*Task) (status string)
	doAllocTiny(size uintptr, task ...*Task) (status string)
	doAllocEnough(size uintptr, task ...*Task) (status string)
	addToList(task *Task)
}

// TODO: remove the manual assignment
func (m *masa) init() *masa {
	m = &masa{nil, 0, 0, MemoryInitializing, nil, &List{nil, Tiny, 0}, nil, nil}
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

func (m *masa) allocate(size uintptr, task ...*Task) {
	_, m.v = m.sliceChecking(task)
	m.task = (m.v).(*Task)

	switch m.level(size) {
	case Little_:
		m.doAllocLittle(size, m.task)
	case Tiny_:
		m.doAllocTiny(size, m.task)
	case Enough_:
		m.doAllocEnough(size, m.task)
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

func (m *masa) sliceChecking(s ...interface{}) (ok bool, v interface{}) {
	if len(s) > 0 {
		m.v = s[0]
		return true, m.v
	}
	return false, nil
}

func (m *masa) checkAllocationStatus(size uintptr, task ...*Task) string {
	if m.isOutOfMemory(size) {
		return MemoryAllocateFailed
	}

	// if current allocation is processing, then we need add
	// task into cached list to wait, if cached list is full,
	// memory allocated failed.
	sign := 0
	if m.cached.ptr != nil {
		sign++
		if m.cached.size < sign {
			return MemoryAllocateFailed
		}
		_, m.v = m.sliceChecking(task)
		m.task = (m.v).(*Task)
		m.addToList(m.task)
		log.Warnf("task cannot add to wait list.")
		return FullOfWaitList
	}
	return MemoryWaitToAllocate
}

func (m *masa) check(size uintptr, task ...*Task) {
	_, m.v = m.sliceChecking(task)
	m.task = (m.v).(*Task)
	m.checkAllocationStatus(size, m.task)
}

// TODO: allocate memory for Tiny level
func (m *masa) doAllocLittle(size uintptr, task ...*Task) (status string) {
	m.check(size, task...)

	return
}

// TODO: allocate memory for Little level
func (m *masa) doAllocTiny(size uintptr, task ...*Task) (status string) {
	m.check(size, task...)

	return
}

// TODO: allocate memory for Enough level
func (m *masa) doAllocEnough(size uintptr, task ...*Task) (status string) {
	m.check(size, task...)

	return
}

// TODO
func (m *masa) addToList(task *Task) {

}
