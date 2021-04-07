// Log the panic under unix to the log file

//+build unix
package log

import (
	"context"
	"os"
	"syscall"
)

func logUncaughtPanic(f *os.File, logger Logger) {
	err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		logger.Panicf(context.Background(), "panic %v", err)
	}
}
