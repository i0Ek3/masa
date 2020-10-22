package masa

import (
	"unsafe"

	log "github.com/sirupsen/logrus"
)

// memory allocation status and size limitations
const (
	MemoryInitializing   = "init"
	MemoryAllocated      = "allocated"
	MemoryAllocating     = "allocating"
	MemoryAllocateFailed = "failed"
	MemoryWaitToAllocate = "waiting"

	FullOfWaitList = "full"

	Tiny   = 16
	Little = 32

	T = iota + 1
	L
	E
)

type emptymem uintptr

// Task defines the task instance
type Task uintptr

// List defines a list
type List struct {
	ptr      unsafe.Pointer
	size     int
	priority int
}

type limbo struct {
	_p *p
}

type p struct {
	_P unsafe.Pointer
}

// masa defines a memory instance
type masa struct {
	ptr       *limbo
	offset    uint
	size      int
	status    string
	allocated [3]bool

	mList  []*Task
	cached *List
	task   *Task
	v      interface{}
}

// Masa defines masa m interface
type Masa interface {
	init() *masa
	getAllocationSize() int
	setAllocationSize(value int) int
	level() int
	allocate(size int, task ...*Task)
	sliceChecking(s interface{}) (bool, interface{})
	checkAllocationStatus(size int, task ...*Task) string
	isOutOfMemory(size int) bool
	check(size int, task ...*Task)
	doAllocLittle(size int, task ...*Task) (status string)
	doAllocTiny(size int, task ...*Task) (status string)
	doAllocEnough(size int, task ...*Task) (status string)
	addToList(task *Task) bool
}

func (m *masa) init() *masa {
	m.ptr = nil
	m.offset = 0
	m.size = 0
	m.status = MemoryInitializing
	for i := 0; i < 3; i++ {
		m.allocated[i] = false
	}
	m.mList = nil
	m.cached = &List{
		ptr:      nil,
		size:     Tiny,
		priority: 0,
	}
	m.task = nil
	m.v = nil
	return m
}

func (m *masa) getAllocationSize() int {
	return m.size
}

func (m *masa) setAllocationSize(value int) int {
	m.size = value
	return m.size
}

func (m *masa) level(size int) int {
	if size <= Little {
		if size >= Tiny {
			return L
		}
		return T
	}
	return E
}

func (m *masa) allocate(size int, task ...*Task) {
	_, m.v = m.sliceChecking(task)
	m.task = (m.v).(*Task)

	switch m.level(size) {
	case L:
		m.doAllocLittle(size, m.task)
	case T:
		m.doAllocTiny(size, m.task)
	case E:
		m.doAllocEnough(size, m.task)
	}
}

func (m *masa) isOutOfMemory(size int) bool {
	if size < Tiny {
		cnt := 0
		for i := 0; i < size; i++ {
			if m.mList[i] != nil {
				cnt++
			}
		}
		if cnt < size {
			return true
		}
		return false
	} else if size >= Tiny && size <= Little {
		cnt := 0
		for i := Tiny; i <= Little; i++ {
			if m.mList[i] != nil {
				cnt++
			}
		}
		if cnt < size {
			return true
		}
		return false
	} else {
		var i int
		cnt := 0
		for ; i > Little; i++ {
			if m.mList[i] != nil {
				cnt++
			}
		}
		if cnt < size {
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

func (m *masa) checkAllocationStatus(size int, task ...*Task) string {
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
		if !m.addToList(m.task) {
			log.Warnf("task cannot add to wait list.")
			return FullOfWaitList
		}
	}
	return MemoryWaitToAllocate
}

func (m *masa) check(size int, task ...*Task) (status string) {
	_, m.v = m.sliceChecking(task)
	m.task = (m.v).(*Task)
	return m.checkAllocationStatus(size, m.task)
}

// allocate memory for Tiny level
func (m *masa) doAllocTiny(size int, task ...*Task) (status string) {
	status = m.check(size, task...)
	if m.ptr != nil && 0 < size && size < Tiny {
		for i := 0; i < size; i++ {
			m.allocated[i] = true
		}
	}
	return
}

// allocate memory for Little level
func (m *masa) doAllocLittle(size int, task ...*Task) (status string) {
	status = m.check(size, task...)
	if m.ptr._p != nil && Tiny <= size && size <= Little {
		for i := Tiny; i < size; i++ {
			m.allocated[i] = true
		}
	}
	return
}

// allocate memory for Enough level
func (m *masa) doAllocEnough(size int, task ...*Task) (status string) {
	status = m.check(size, task...)
	if m.ptr._p._P != nil && size > Little {
		for i := Little; ; i++ {
			m.allocated[i] = true
		}
	}
	return
}

func (m *masa) addToList(task *Task) bool {
	if task != nil {
		m.mList = append(m.mList, task)
		return true
	}
	return false
}
