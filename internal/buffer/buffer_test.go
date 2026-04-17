package buffer

import (
	"fmt"
	"testing"
)

func TestBuffer_WriteAndFlush(t *testing.T) {
	b := New(4)
	b.Write("a")
	b.Write("b")
	b.Write("c")

	lines := b.Flush()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "a" || lines[1] != "b" || lines[2] != "c" {
		t.Errorf("unexpected lines: %v", lines)
	}
	if b.Len() != 0 {
		t.Errorf("expected empty buffer after flush")
	}
}

func TestBuffer_Overflow_DropsOldest(t *testing.T) {
	b := New(3)
	for i := 0; i < 5; i++ {
		b.Write(fmt.Sprintf("line%d", i))
	}

	if b.Dropped() != 2 {
		t.Errorf("expected 2 dropped, got %d", b.Dropped())
	}

	lines := b.Flush()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	// oldest two (line0, line1) should be dropped
	if lines[0] != "line2" {
		t.Errorf("expected line2, got %s", lines[0])
	}
}

func TestBuffer_FlushEmpty(t *testing.T) {
	b := New(8)
	lines := b.Flush()
	if len(lines) != 0 {
		t.Errorf("expected empty slice, got %v", lines)
	}
}

func TestBuffer_ZeroCapacity_DefaultsToOne(t *testing.T) {
	b := New(0)
	b.Write("only")
	lines := b.Flush()
	if len(lines) != 1 || lines[0] != "only" {
		t.Errorf("unexpected result: %v", lines)
	}
}

func TestBuffer_Len(t *testing.T) {
	b := New(10)
	if b.Len() != 0 {
		t.Errorf("expected 0")
	}
	b.Write("x")
	b.Write("y")
	if b.Len() != 2 {
		t.Errorf("expected 2, got %d", b.Len())
	}
}

func TestBuffer_NoDrop_NoDropCount(t *testing.T) {
	b := New(5)
	b.Write("a")
	b.Write("b")
	if b.Dropped() != 0 {
		t.Errorf("expected 0 dropped")
	}
}
