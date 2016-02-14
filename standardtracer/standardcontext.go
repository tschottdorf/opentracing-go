package standardtracer

import "sync"

type locker struct {
	sync.RWMutex
	traceAttrs map[string]string // initialized on first use
}

// StandardContext holds the basic Span metadata.
type StandardContext struct {
	// A probabilistically unique identifier for a [multi-span] trace.
	TraceID int64

	// A probabilistically unique identifier for a span.
	SpanID int64

	// The SpanID of this StandardContext's parent, or 0 if there is no parent.
	ParentSpanID int64

	// Whether the trace is sampled.
	Sampled bool

	// `tagLock` protects the `traceAttrs` map, which in turn supports
	// `SetTraceAttribute` and `TraceAttribute`.
	mu locker
}

// NewRootStandardContext creates a StandardContext corresponding to a root
// span.
func NewRootStandardContext() StandardContext {
	return StandardContext{
		TraceID: randomID(),
		SpanID:  randomID(),
		Sampled: randomID()%64 == 0,
		mu:      locker{traceAttrs: make(map[string]string)},
	}
}

// NewChild creates a new child StandardContext.
func (c *StandardContext) NewChild() StandardContext {
	c.mu.RLock()
	newTags := make(map[string]string, len(c.mu.traceAttrs))
	for k, v := range c.mu.traceAttrs {
		newTags[k] = v
	}
	c.mu.RUnlock()

	return StandardContext{
		TraceID:      c.TraceID,
		SpanID:       randomID(),
		ParentSpanID: c.SpanID,
		Sampled:      c.Sampled,
		mu:           locker{traceAttrs: newTags},
	}
}
