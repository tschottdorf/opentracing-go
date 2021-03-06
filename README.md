# OpenTracing API for Go

This package is a Go platform API for OpenTracing.

## Required Reading

In order to understand the Go platform API, one must first be familiar with the
[OpenTracing project](http://opentracing.io) and
[terminology](http://opentracing.io/spec/) more generally.

## API overview for those adding instrumentation

Everyday consumers of this `opentracing` package really only need to worry
about a couple of key abstractions: the `StartTrace` function, the `Span`
interface, and binding a `Tracer` at `main()`-time. Here are code snippets
demonstrating some important use cases.

#### Singleton initialization

The simplest starting point is `./default_tracer.go`. As early as possible, call

```go
    import ".../opentracing-go"
    import ".../some_tracing_impl"
    
    func main() {
        tracerImpl := some_tracing_impl.New(...) // tracing impl specific
        opentracing.InitGlobalTracer(tracerImpl)
        ...
    }
```

##### Non-Singleton initialization

If you prefer direct control to singletons, manage ownership of the
`opentracing.Tracer` implementation explicitly.

#### Starting an empty trace by creating a "root span"

```go
    func xyz() {
        ...
        sp := opentracing.StartTrace("span_name")
        defer sp.Finish()
        sp.LogEvent("xyz_called")
        ...
    }
```

#### Creating a Span given an existing Span

```go
    func xyz(parentSpan opentracing.Span, ...) {
        ...
        sp := opentracing.JoinTrace("span_name", parentSpan)
        defer sp.Finish()
        sp.LogEvent("xyz_called")
        ...
    }
```

#### Creating a Span given an existing Golang `context.Context`

Additionally, this example demonstrates how to get a `context.Context`
associated with any `opentracing.Span` instance.

```go
    func xyz(goCtx context.Context, ...) {
        ...
        goCtx, sp := opentracing.ContextWithSpan(
            goCtx, opentracing.JoinTrace("span_name", goCtx))
        defer sp.Finish()
        sp.LogEvent("xyz_called")
        ...
    }
```

#### Serializing to the wire

```go
    func makeSomeRequest(ctx context.Context) ... {
        if span := opentracing.SpanFromContext(ctx); span != nil {
            httpClient := &http.Client{}
            httpReq, _ := http.NewRequest("GET", "http://myservice/", nil)

            // Transmit the span's TraceContext as an HTTP header on our
            // outbound request.
            opentracing.GlobalTracer().PropagateSpanInHeader(
                span,
                httpReq.Header,
                opentracing.DefaultTracer())

            resp, err := httpClient.Do(httpReq)
            ...
        }
        ...
    }
```

#### Deserializing from the wire

```go
    http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        // Join the trace in the HTTP header using the opentracing helper.
        serverSpan, err := opentracing.GlobalTracer().JoinTraceFromHeader(
                "serverSpan", req.Header, opentracing.GlobalTracer())
        if err != nil {
            serverSpan = opentracing.StartTrace("serverSpan")
        }
        var goCtx context.Context = ...
        goCtx, _ = opentracing.ContextWithSpan(goCtx, serverSpan)
        defer serverSpan.Finish()
        ...
    }
```

#### Goroutine-safety

The entire public API is goroutine-safe and does not require external
synchronization.

## API pointers for those implementing a tracing system

Tracing system implementors may be able to reuse or copy-paste-modify the `./standardtracer` package. In particular, see `standardtracer.New(...)`.
