package trace

import (
	"io"
	"fmt"
)

// Tracer is the interface that describes an object capable of
// tracing events throughout code.
type Tracer interface {
	Trace(...interface{})
}

// New is a helper function for creating an new Tracer instance.
// It use the built-in private tracer object so that its method
// will be exposed only.
func New(writer io.Writer) Tracer {
	return &tracer{
		writer: writer,
	}
}

// Off is an helper function that return an implementation
// Tracer that print nothing.
func Off() Tracer {
	return nilTracer{}
}

// nilTracer is an implementation of the Tracer interface
// use for turning off the tracing.
type nilTracer struct {}

// Trace will do nothing.
func (nilTracer) Trace(args ...interface{}) {}

// tracer is a private object that will be use during
// the initialization of the new Tracer interface
// it exposes only its method.
type tracer struct {
	// writer is an interface where the trace logs will
	// be written.
	writer io.Writer
}

// Trace satisfy the Tracer interface.
func (t *tracer) Trace(args ...interface{}) {
	fmt.Fprint(t.writer, args...)
	fmt.Fprintln(t.writer)
}