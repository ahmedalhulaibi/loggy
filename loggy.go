package loggy

import (
	"context"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
	fields []string
}

func New(zapLogger *zap.SugaredLogger) Logger {
	return Logger{
		SugaredLogger: zapLogger,
	}
}

// WithFields returns a new logger which will inspect for additional fields
func (l Logger) WithFields(fields ...string) Logger {
	l.fields = append(l.fields, fields...)
	return l
}

type KeyVal struct {
	Key string
	Val interface{}
}

func (l Logger) Log(ctx context.Context, msg string, args ...KeyVal) {
	argsI := make([]interface{}, 0, len(args)*2)

	for _, kv := range args {
		argsI = append(argsI, kv.Key, kv.Val)
	}

	finalArgs := append(l.extractArgs(ctx), argsI...)
	l.Infow(msg, finalArgs...)
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
	return l.SugaredLogger.Sync()
}
