package masa

import (
	"testing"
)

func TestDoAllocate(t *testing.T) {
	check := func(t *testing.T, got, want string) {
		if got != want {
			t.Errorf("got '%q' want '%q'", got, want)
		}
	}

	var (
		m    masa
		task *Task
	)

	t.Run("doAllocTiny", func(t *testing.T) {
		got := m.doAllocTiny(10, task)
		want := MemoryAllocated
		check(t, got, want)
	})

	t.Run("doAllocTiny", func(t *testing.T) {
		got := m.doAllocTiny(1024, task)
		want := MemoryAllocateFailed
		check(t, got, want)
	})

	t.Run("doAllocLittle", func(t *testing.T) {
		got := m.doAllocLittle(18, task)
		want := MemoryAllocated
		check(t, got, want)
	})

	t.Run("doAllocLittle", func(t *testing.T) {
		got := m.doAllocLittle(1024, task)
		want := MemoryAllocateFailed
		check(t, got, want)
	})

	t.Run("doAllocEnough", func(t *testing.T) {
		got := m.doAllocEnough(1024, task)
		want := MemoryAllocated
		check(t, got, want)
	})

	t.Run("doAllocEnough", func(t *testing.T) {
		got := m.doAllocEnough(20, task)
		want := MemoryAllocateFailed
		check(t, got, want)
	})
}
