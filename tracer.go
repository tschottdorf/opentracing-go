package opentracing

// Tracer is a simple, thin interface for Span creation.
//
// A straightforward implementation is available via the
// `opentracing/standardtracer` package's `standardtracer.New()'.
type Tracer interface {
	SpanPropagator

	// Create, start, and return a new Span with the given `operationName`, all
	// without specifying a parent Span that can be used to incorporate the
	// newly-returned Span into an existing trace. (I.e., the returned Span is
	// the "root" of its trace).
	//
	// Examples:
	//
	//     var tracer Tracer = ...
	//
	//     sp := tracer.StartTrace("GetFeed")
	//
	//     sp := tracer.StartTrace("HandleHTTPRequest").
	//         SetTag("user_agent", req.UserAgent).
	//         SetTag("lucky_number", 42)
	//
	StartTrace(operationName string) Span
}
