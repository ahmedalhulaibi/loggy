package loggy

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
	fields           []string
	logLevelSelector map[zapcore.Level]func(msg string, keysAndValues ...interface{})
}

func New(zapLogger *zap.SugaredLogger) Logger {
	logLevelSelector := map[zapcore.Level]func(msg string, keysAndValues ...interface{}){
		zap.DebugLevel:  zapLogger.Debugw,
		zap.InfoLevel:   zapLogger.Infow,
		zap.WarnLevel:   zapLogger.Warnw,
		zap.ErrorLevel:  zapLogger.Errorw,
		zap.DPanicLevel: zapLogger.DPanicw,
		zap.PanicLevel:  zapLogger.Panicw,
		zap.FatalLevel:  zapLogger.Fatalw,
	}

	return Logger{
		SugaredLogger:    zapLogger,
		logLevelSelector: logLevelSelector,
	}
}

// WithFields returns a new logger which will inspect for additional fields
func (l Logger) WithFields(fields ...string) Logger {
	l.fields = append(l.fields, fields...)
	return l
}

func (l Logger) WithLogLevelSelector(s map[zapcore.Level]func(string, ...interface{})) Logger {
	if s != nil {
		l.logLevelSelector = s
	}
	return l
}

type KeyVal struct {
	Key string
	Val interface{}
}

// TODO: Instead wrap the zap.SugaredLogger methods, this is redundant
func (l Logger) Log(ctx context.Context, level zapcore.Level, msg string, args ...KeyVal) {
	argsI := make([]interface{}, 0, len(args)*2)

	for _, kv := range args {
		argsI = append(argsI, kv.Key, kv.Val)
	}

	finalArgs := append(l.extractArgs(ctx), argsI...)

	if logger, ok := l.logLevelSelector[level]; ok {
		logger(msg, finalArgs...)
	}
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
