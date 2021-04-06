package log

import (
	"context"
	"fmt"
)

type printLogger struct {
	withLogName bool
	format      string
}

//PrintLogger print logger
func PrintLogger(ln bool) Logger {
	format := printFormatter
	if ln {
		format += "|%v"
	}
	return &printLogger{withLogName: ln, format: format}
}

//Info info
func (l printLogger) Info(ctx context.Context, format string, args ...interface{}) {
	_log(ctx, l, "INFO", format, args...)
}

//Error error
func (l printLogger) Error(ctx context.Context, format string, args ...interface{}) {
	_log(ctx, l, "ERROR", format, args...)
}

//Warning warning
func (l printLogger) Warning(ctx context.Context, format string, args ...interface{}) {
	_log(ctx, l, "WARN", format, args...)
}

//Debug debug
func (l printLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	_log(ctx, l, "DEBUG", format, args...)
}

//Panicf panics
func (l printLogger) Panicf(ctx context.Context, format string, arg ...interface{}) {
	l.Error(ctx, format, arg...)
	panic(fmt.Sprintf(format, arg...))
}
