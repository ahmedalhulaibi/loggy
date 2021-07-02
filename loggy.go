package loggy

import (
	"context"

	"go.uber.org/zap"
)

type Logger struct {
	zapLogger *zap.SugaredLogger
	fields    []string
}

func New(zapLogger *zap.SugaredLogger) Logger {
	return Logger{
		zapLogger: zapLogger,
	}
}

// WithFields returns a new logger which will inspect for additional fields
func (l Logger) WithFields(fields ...string) Logger {
	l.fields = append(l.fields, fields...)
	return l
}

func (l Logger) Log(ctx context.Context, msg string, args ...interface{}) {
	finalArgs := append(l.extractArgs(ctx), args...)
	l.zapLogger.Infow(msg, finalArgs...)
}

func (l Logger) extractArgs(ctx context.Context) []interface{} {
	var ctxArgs []interface{}
	for _, field := range l.fields {
		val := ctx.Value(field)
		if val != nil {
			ctxArgs = append(ctxArgs, field, val)
		}
	}
	return ctxArgs
}

func (l Logger) Sync() error {
	return l.zapLogger.Sync()
}
