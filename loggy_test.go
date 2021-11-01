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
		givenFields  []LoggyCtxKey
		givenCtx     context.Context
		expectedArgs []interface{}
	}{
		"Should return empty slice": {
			givenCtx: context.Background(),
		},
		"Should return request_id field and value": {
			givenFields:  []LoggyCtxKey{"request_id"},
			givenCtx:     context.WithValue(context.Background(), "request_id", "<request-id-value>"),
			expectedArgs: []interface{}{"request_id", "<request-id-value>"},
		},
		"Should return not return request_id field and value if not present": {
			givenFields: []LoggyCtxKey{"request_id"},
			givenCtx:    context.Background(),
		},
		"Should return request_id field and value and user field and value": {
			givenFields:  []LoggyCtxKey{"request_id", "user"},
			givenCtx:     context.WithValue(context.WithValue(context.Background(), "request_id", "<request-id-value>"), "user", "<user-value>"),
			expectedArgs: []interface{}{"request_id", "<request-id-value>", "user", "<user-value>"},
		},
		"Should only return user field and value if request_id not present": {
			givenFields:  []LoggyCtxKey{"request_id", "user"},
			givenCtx:     context.WithValue(context.Background(), "user", "<user-value>"),
			expectedArgs: []interface{}{"user", "<user-value>"},
		},
	}

	for n, tc := range tests {
		tc := tc
		t.Run(n, func(t *testing.T) {
			zaplogger, err := zap.NewProduction()
			require.NoError(t, err)

			sugaredLogger := zaplogger.Sugar()

			logger := New(sugaredLogger)
			logger = logger.WithFields(tc.givenFields...)

			actualArgs := logger.extractArgsFromCtx(tc.givenCtx)
			require.ElementsMatch(t, tc.expectedArgs, actualArgs)
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	nopLogger := zap.NewNop().Sugar()
	type fields struct {
		SugaredLogger *zap.SugaredLogger
		fields        []LoggyCtxKey
	}
	type args struct {
		fields []LoggyCtxKey
	}
	tests := map[string]struct {
		fields fields
		args   args
		want   Logger
	}{
		"Should return logger with new fields appended": {
			fields: fields{
				fields:        []LoggyCtxKey{"request_id", "user_id"},
				SugaredLogger: nopLogger,
			},
			args: args{fields: []LoggyCtxKey{"run_id"}},
			want: Logger{
				fields:        []LoggyCtxKey{"request_id", "user_id", "run_id"},
				SugaredLogger: nopLogger,
			},
		},
		"Should return logger with no new fields appended": {
			fields: fields{
				fields:        []LoggyCtxKey{"request_id", "user_id"},
				SugaredLogger: nopLogger,
			},
			args: args{fields: []LoggyCtxKey{}},
			want: Logger{
				fields:        []LoggyCtxKey{"request_id", "user_id"},
				SugaredLogger: nopLogger,
			},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			l := Logger{
				SugaredLogger: tc.fields.SugaredLogger,
				fields:        tc.fields.fields,
			}
			got := l.WithFields(tc.args.fields...)
			require.Equal(t, tc.want, got)
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
			l := New(zapLogger.Sugar()).WithFields("request_id")

			ctx := context.WithValue(context.Background(), "request_id", "<request-id-value>")
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

			l := New(zapLogger.Sugar()).WithFields("request_id")

			ctx := context.WithValue(context.Background(), "request_id", "<request-id-value>")

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
		logFunc func(Logger, context.Context, string, ...KeyVal)
	}{
		"Should log with debug level": {
			logFunc: func(l Logger, ctx context.Context, msg string, args ...KeyVal) {
				l.Debugw(ctx, msg, args...)
			},
		},
		"Should log with info level": {
			logFunc: func(l Logger, ctx context.Context, msg string, args ...KeyVal) {
				l.Infow(ctx, msg, args...)
			},
		},
		"Should log with warn level": {
			logFunc: func(l Logger, ctx context.Context, msg string, args ...KeyVal) {
				l.Warnw(ctx, msg, args...)
			},
		},
		"Should log with error level": {
			logFunc: func(l Logger, ctx context.Context, msg string, args ...KeyVal) {
				l.Errorw(ctx, msg, args...)
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})

			zapLogger := newZapTestLogger(t, zapcore.AddSync(buf))

			l := New(zapLogger.Sugar()).WithFields("request_id")

			ctx := context.WithValue(context.Background(), "request_id", "<request-id-value>")

			tc.logFunc(l, ctx, "something goes here", KeyVal{Key: "key", Val: "value"})

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
