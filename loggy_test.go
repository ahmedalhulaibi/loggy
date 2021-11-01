package loggy

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestExtractArgs(t *testing.T) {
	tests := map[string]struct {
		givenFields  []interface{}
		givenCtx     context.Context
		expectedArgs []interface{}
	}{
		"Should return empty slice": {
			givenCtx: context.Background(),
		},
		"Should return request_id field and value": {
			givenCtx:     ContextWithArgs(context.Background(), "request_id", "<request-id-value>"),
			expectedArgs: []interface{}{"request_id", "<request-id-value>"},
		},
		"Should return not return request_id field and value if not present": {
			givenCtx: context.Background(),
		},
		"Should return request_id field and value and user field and value": {
			givenCtx:     ContextWithArgs(context.Background(), "request_id", "<request-id-value>", "user", "<user-value>"),
			expectedArgs: []interface{}{"request_id", "<request-id-value>", "user", "<user-value>"},
		},
		"Should only return user field and value if request_id not present": {
			givenCtx:     ContextWithArgs(context.Background(), "user", "<user-value>"),
			expectedArgs: []interface{}{"user", "<user-value>"},
		},
	}

	for n, tc := range tests {
		tc := tc
		t.Run(n, func(t *testing.T) {
			actualArgs := extractArgsFromCtx(tc.givenCtx)
			require.ElementsMatch(t, tc.expectedArgs, actualArgs)
		})
	}
}

func TestLogger_LogUntemplatedMessage(t *testing.T) {
	tests := map[string]struct {
		logFunc func(Logger, context.Context, ...interface{})
	}{
		"Should log with debug level": {
			logFunc: func(l Logger, ctx context.Context, args ...interface{}) {
				l.Debug(ctx, args...)
			},
		},
		"Should log with info level": {
			logFunc: func(l Logger, ctx context.Context, args ...interface{}) {
				l.Info(ctx, args...)
			},
		},
		"Should log with warn level": {
			logFunc: func(l Logger, ctx context.Context, args ...interface{}) {
				l.Warn(ctx, args...)
			},
		},
		"Should log with error level": {
			logFunc: func(l Logger, ctx context.Context, args ...interface{}) {
				l.Error(ctx, args...)
			},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})

			zapLogger := newZapTestLogger(t, zapcore.AddSync(buf))
			l := New(zapLogger.Sugar())

			ctx := ContextWithArgs(context.Background(), "request_id", "<request-id-value>")
			tc.logFunc(l, ctx, "something goes here")

			if *updateGolden {
				t.Log("Updating golden file:", goldenFilename(t))
				require.NoError(t, os.MkdirAll(filepath.Dir(goldenFilename(t)), 0755))
				require.NoError(t, os.WriteFile(goldenFilename(t), buf.Bytes(), 0644))
			}

			golden, err := os.ReadFile(goldenFilename(t))
			require.NoError(t, err)
			require.Equal(t, buf.Bytes(), golden)
		})
	}
}

func TestLogger_LogTemplatedMessage(t *testing.T) {
	tests := map[string]struct {
		logFunc func(Logger, context.Context, string, ...interface{})
	}{
		"Should log with debug level": {
			logFunc: func(l Logger, ctx context.Context, template string, args ...interface{}) {
				l.Debugf(ctx, template, args...)
			},
		},
		"Should log with info level": {
			logFunc: func(l Logger, ctx context.Context, template string, args ...interface{}) {
				l.Infof(ctx, template, args...)
			},
		},
		"Should log with warn level": {
			logFunc: func(l Logger, ctx context.Context, template string, args ...interface{}) {
				l.Warnf(ctx, template, args...)
			},
		},
		"Should log with error level": {
			logFunc: func(l Logger, ctx context.Context, template string, args ...interface{}) {
				l.Errorf(ctx, template, args...)
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})

			zapLogger := newZapTestLogger(t, zapcore.AddSync(buf))

			l := New(zapLogger.Sugar())

			ctx := ContextWithArgs(context.Background(), "request_id", "<request-id-value>")

			tc.logFunc(l, ctx, "something goes here %s", "here")

			if *updateGolden {
				t.Log("Updating golden file:", goldenFilename(t))
				require.NoError(t, os.MkdirAll(filepath.Dir(goldenFilename(t)), 0755))
				require.NoError(t, os.WriteFile(goldenFilename(t), buf.Bytes(), 0644))
			}

			golden, err := os.ReadFile(goldenFilename(t))
			require.NoError(t, err)
			require.Equal(t, buf.Bytes(), golden)
		})
	}
}

func TestLogger_LogMessageWithFields(t *testing.T) {
	tests := map[string]struct {
		logFunc func(Logger, context.Context, string, ...interface{})
	}{
		"Should log with debug level": {
			logFunc: func(l Logger, ctx context.Context, msg string, args ...interface{}) {
				l.Debugw(ctx, msg, args...)
			},
		},
		"Should log with info level": {
			logFunc: func(l Logger, ctx context.Context, msg string, args ...interface{}) {
				l.Infow(ctx, msg, args...)
			},
		},
		"Should log with warn level": {
			logFunc: func(l Logger, ctx context.Context, msg string, args ...interface{}) {
				l.Warnw(ctx, msg, args...)
			},
		},
		"Should log with error level": {
			logFunc: func(l Logger, ctx context.Context, msg string, args ...interface{}) {
				l.Errorw(ctx, msg, args...)
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})

			zapLogger := newZapTestLogger(t, zapcore.AddSync(buf))

			l := New(zapLogger.Sugar())

			ctx := ContextWithArgs(context.Background(), "request_id", "<request-id-value>")

			tc.logFunc(l, ctx, "something goes here", "key", "value")

			if *updateGolden {
				t.Log("Updating golden file:", goldenFilename(t))
				require.NoError(t, os.MkdirAll(filepath.Dir(goldenFilename(t)), 0755))
				require.NoError(t, os.WriteFile(goldenFilename(t), buf.Bytes(), 0644))
			}

			golden, err := os.ReadFile(goldenFilename(t))
			require.NoError(t, err)
			require.Equal(t, buf.Bytes(), golden)
		})
	}
}

func goldenFilename(t *testing.T) string {
	t.Helper()
	return "testdata/" + t.Name() + ".golden"
}

func newZapTestLogger(t *testing.T, output zapcore.WriteSyncer, options ...zap.Option) *zap.Logger {
	t.Helper()
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), output, zap.DebugLevel)
	return zap.New(core).WithOptions(options...)
}

// BenchmarkLoggy benchmarks the recommended usage of the Logger.
// It is intended to be run with the -benchmem flag.
// The recommended usage of the Logger is to use the WithFields and Infow, Debugw, etc. methods.
func BenchmarkLoggy(b *testing.B) {
	// The Logger allocation is not included in the benchmark time since it is declared once at the beginning of the program
	// It is expected that in the real world the Logger will be allocated once and reused across the application.
	l := New(zap.NewNop().Sugar())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// It is expected that context would still be modified in middleware with request-scoped values
		ctx := ContextWithArgs(context.Background(), "request_id", "<request-id-value>")

		// Elsewhere in the codebase, the same instance of logger can be used and will extract request-scoped values from context.Context
		// For the sake of the test, let's assume that we log ten times per request.
		l.Infow(ctx, "something goes here", "key", "value")
		l.Infow(ctx, "something goes here", "key", "value")
		l.Infow(ctx, "something goes here", "key", "value")
		l.Infow(ctx, "something goes here", "key", "value")
		l.Infow(ctx, "something goes here", "key", "value")
		l.Infow(ctx, "something goes here", "key", "value")
		l.Infow(ctx, "something goes here", "key", "value")
		l.Infow(ctx, "something goes here", "key", "value")
		l.Infow(ctx, "something goes here", "key", "value")
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
		// For the sake of the test, let's assume that we log ten times per request.
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
		if maybeLogger, ok := ctx.Value("logger").(*zap.SugaredLogger); ok {
			maybeLogger.Infow("something goes here", "key", "value")
		}
	}
}
