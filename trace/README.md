# Trace

This package provides an interface for recording the latency of operations and logging details
about all operations where the latency exceeds a limit.

## Usage

To create a trace:

```go
func doSomething() {
    ctx, opTrace := trace.New(context.Background(), "operation", attribute.String("fieldKey1", "fieldValue1"))
    defer opTrace.LogIfLong(ctx, 100 * time.Millisecond)
    // do something
}
```

To split an trace into multiple steps:

```go
func doSomething() {
    ctx, opTrace := trace.New(context.Background(), "operation")
    defer opTrace.LogIfLong(ctx, 100 * time.Millisecond)
    // do step 1
    opTrace.Step(ctx, "step1", Field{Key: "stepFieldKey1", Value: "stepFieldValue1"})
    // do step 2
    opTrace.Step(ctx, "step2")
}
```

To nest traces:

```go
func doSomething() {
    ctx, rootTrace := trace.New(context.Background(), "rootOperation")
    defer rootTrace.LogIfLong(ctx, 100 * time.Millisecond)
    
    func() {
        ctx, nestedTrace := rootTrace.Nest(ctx, "nested", Field{Key: "nestedFieldKey1", Value: "nestedFieldValue1"})
        defer nestedTrace.LogIfLong(ctx, 50 * time.Millisecond)
        // do nested operation
    }()
}
```

Traces can also be logged unconditionally or introspected:

```go
opTrace.TotalTime() // Duration since the Trace was created
opTrace.Log(context.Background()) // unconditionally log the trace
```

### Using context.Context to nest traces

`context.Context` can be used to manage nested traces. Create traces by calling `trace.GetTraceFromContext(ctx).Nest`. 
This is safe even if there is no parent trace already in the context because `(*(Trace)nil).Nest()` returns
a top level trace.

```go
func doSomething(ctx context.Context) {
    ctx, opTrace := trace.FromContext(ctx).Nest(ctx, "operation") // create a trace, possibly nested
    ctx = trace.ContextWithTrace(ctx, opTrace) // make this trace the parent trace of the context
    defer opTrace.LogIfLong(ctx, 50 * time.Millisecond)
    
    doSomethingElse(ctx)
}
```