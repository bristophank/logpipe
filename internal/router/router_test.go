package router

import (
	"bytes"
	"testing"
)

func TestRoute_AllSinks(t *testing.T) {
	r := New()
	var buf1, buf2 bytes.Buffer
	r.AddSink("a", &buf1)
	r.AddSink("b", &buf2)

	if err := r.Route([]byte(`{"level":"info"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, buf := range []*bytes.Buffer{&buf1, &buf2} {
		if buf.Len() == 0 {
			t.Error("expected data written to sink")
		}
	}
}

func TestRoute_NamedSink(t *testing.T) {
	r := New()
	var buf1, buf2 bytes.Buffer
	r.AddSink("a", &buf1)
	r.AddSink("b", &buf2)

	if err := r.Route([]byte(`{"level":"error"}`), "a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf1.Len() == 0 {
		t.Error("expected data written to sink a")
	}
	if buf2.Len() != 0 {
		t.Error("expected sink b to be empty")
	}
}

func TestRoute_MissingSink(t *testing.T) {
	r := New()
	var buf bytes.Buffer
	r.AddSink("a", &buf)

	if err := r.Route([]byte(`{"level":"warn"}`), "nonexistent"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Error("expected no data written")
	}
}

func TestRoute_RemoveSink(t *testing.T) {
	r := New()
	var buf bytes.Buffer
	r.AddSink("a", &buf)
	r.RemoveSink("a")

	if err := r.Route([]byte(`{"level":"debug"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Error("expected no data after removal")
	}
}

func TestRoute_EmptyLine(t *testing.T) {
	r := New()
	var buf bytes.Buffer
	r.AddSink("a", &buf)

	if err := r.Route([]byte{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Error("expected nothing written for empty line")
	}
}
