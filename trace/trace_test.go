package trace

import (
	"testing"
	"bytes"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Fatal("tracer should not be nil")
	}

	tracer.Trace("Hello trace package")
	if buf.String() != "Hello trace package\n" {
		t.Error("expected message wasn't match")
	}
}
