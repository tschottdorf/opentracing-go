package opentracing

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	// ContextIDHTTPHeaderPrefix precedes the opentracing-related ContextID HTTP
	// headers.
	ContextIDHTTPHeaderPrefix = "Open-Tracing-Context-Id-"

	// TagsHTTPHeaderPrefix precedes the opentracing-related trace-tags HTTP
	// headers.
	TagsHTTPHeaderPrefix = "Open-Tracing-Trace-Tags-"
)

// PropagateSpanInHeader encodes Span `ctx` to `h` as a series of
// HTTP headers. Values are URL-escaped.
func PropagateSpanInHeader(
	sp Span,
	h http.Header,
	propagator SpanPropagator,
) {
	contextIDMap, tagsMap := propagator.PropagateSpanAsText(sp)
	for headerSuffix, val := range contextIDMap {
		h.Add(ContextIDHTTPHeaderPrefix+headerSuffix, url.QueryEscape(val))
	}
	for headerSuffix, val := range tagsMap {
		h.Add(TagsHTTPHeaderPrefix+headerSuffix, url.QueryEscape(val))
	}
}

// NewSpanFromHeader decodes a Span with operation name `operationName` from
// `h`, expecting that header values are URL-escpaed.
//
// If `operationName` is empty, the caller must later call
// `Span.SetOperationName` on the returned `Span`.
func JoinTraceFromHeader(
	operationName string,
	h http.Header,
	propagator SpanPropagator,
) (Span, error) {
	contextIDMap := make(map[string]string)
	tagsMap := make(map[string]string)
	for key, val := range h {
		if strings.HasPrefix(key, ContextIDHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			contextIDMap[strings.TrimPrefix(key, ContextIDHTTPHeaderPrefix)] = unescaped
		} else if strings.HasPrefix(key, TagsHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			tagsMap[strings.TrimPrefix(key, TagsHTTPHeaderPrefix)] = unescaped
		}
	}
	return propagator.JoinTraceFromText(operationName, contextIDMap, tagsMap)
}
