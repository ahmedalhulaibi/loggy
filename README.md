# loggy

`loggy` is an extension of `uber/zap` sugared logger designed to improve ergonomics of logging request-scoped values.

Typically in Go backend codebases, each request (HTTP, gRPC, etc.) has request-scoped values that are logged with each message. The way I have seen this implemented in many codebases is by extended the logging context using a method `log.With` or `log.WithFields` which returns a new logger instance. The new logger is injected into `context.Context`. Example below:

`zap` with `loggy` changes the semantics slightly. Rather than creating a new request-scoped logger, the logger itself is configured to search for specific request-scoped values. `loggy.Logger` accepts `context.Context` and extracts request-scoped values directly from `context.Context`. If a field does not exist in `context.Context` it is ignored.

This solves two issues:
1. No need to inject a service dependency - a logger - via `context.Context` which can lead to panics at runtime if implemented incorrectly
2. No need to allocate a new logger with each request, which can have a marginal improvement in performance

To be clear, the issues above are not an issue with the logging libraries themselves, but an issue with the logging practices established.

# Benchmarks

```
goos: darwin
goarch: amd64
pkg: github.com/ahmedalhulaibi/loggy
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
BenchmarkLoggy-8        98274076               247.4 ns/op           176 B/op          5 allocs/op
BenchmarkZap-8          92253990               269.1 ns/op           280 B/op          4 allocs/op
PASS
ok      github.com/ahmedalhulaibi/loggy 49.858s
```

```go
// BenchmarkLoggy benchmarks the recommended usage of the Logger.
// It is intended to be run with the -benchmem flag.
// The recommended usage of the Logger is to use the WithFields and Infow, Debugw, etc. methods.
func BenchmarkLoggy(b *testing.B) {
	// The Logger allocation is not included in the benchmark time since it is declared once at the beginning of the program
	// It is expected that in the real world the Logger will be allocated once and reused across the application.
	l := New(zap.NewNop().Sugar()).WithFields("request_id")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// It is expected that context would still be modified in middleware with request-scoped values
		ctx := context.WithValue(context.Background(), "request_id", "<request-id-value>")

		// Elsewhere in the codebase, the same instance of logger can be used and will extract request-scoped values from context.Context
		l.Infow(ctx, "something goes here", "key", "value")
	}
}

// BenchmarkZap benchmarks the usage of the zap logger as it would be in the real world.
// It is intended to be run with the -benchmem flag.
func BenchmarkZap(b *testing.B) {
	l := zap.NewNop().Sugar()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Typically the zap logger is injected into context with request-scoped fields in middleware
		ctx := context.WithValue(context.Background(), "logger", l.With("request_id", "<request-id-value>"))

		// Elsewhere in the codebase we can extract and use the specific request-scoped logger
		// Typically this extract logic is wrapped in a helper e.g. logger(ctx).Infow but that is not relevant to this benchmark
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
	}
}
```