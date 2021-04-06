package main

import (
	"context"

	"github.com/RomanLorens/logger/log"
)

func main() {
	logger, _ := log.New(log.WithConfig("test.log").
		WithMaxSize(10).WithMaxBackups(8).Build())
	logger2, _ := log.New(log.WithConfig("test-2.log").WithLogName(true).Build())

	pl := log.PrintLogger(false)
	logger.Info(context.Background(), "boom")
	logger.Debug(context.Background(), "debug %v", 23)
	logger2.Info(context.Background(), "log name format")
	pl.Info(context.Background(), "print logger")
}
