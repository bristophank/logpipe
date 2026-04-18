package tee_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/yourorg/logpipe/internal/tee"
)

func TestWrite_NoWriters(t *testing.T) {
	te := tee.New()
	n, err := te.Write([]byte("hello\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 6 {
		t.Fatalf("expected 6 got %d", n)
	}
}

func TestWrite_MultipleWriters(t *testing.T) {
	te := tee.New()
	var a, b bytes.Buffer
	te.Add("a", &a)
	te.Add("b", &b)

	msg := []byte(`{"level":"info"}` + "\n")
	te.Write(msg)

	if a.String() != string(msg) {
		t.Errorf("writer a: got %q", a.String())
	}
	if b.String() != string(msg) {
		t.Errorf("writer b: got %q", b.String())
	}
}

func TestWrite_RemoveWriter(t *testing.T) {
	te := tee.New()
	var a bytes.Buffer
	te.Add("a", &a)
	te.Remove("a")

	te.Write([]byte("data\n"))
	if a.Len() != 0 {
		t.Error("expected no data after remove")
	}
}

func TestLen(t *testing.T) {
	te := tee.New()
	if te.Len() != 0 {
		t.Fatal("expected 0")
	}
	te.Add("x", io.Discard)
	te.Add("y", io.Discard)
	if te.Len() != 2 {
		t.Fatalf("expected 2 got %d", te.Len())
	}
}

type errWriter struct{ err error }

func (e *errWriter) Write(p []byte) (int, error) { return 0, e.err }

func TestWrite_PropagatesFirstError(t *testing.T) {
	te := tee.New()
	expected := errors.New("sink down")
	te.Add("bad", &errWriter{err: expected})
	_, err := te.Write([]byte("line\n"))
	if err != expected {
		t.Fatalf("expected sink error, got %v", err)
	}
}
