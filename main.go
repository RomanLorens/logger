package main

import (
	"context"

	"github.com/RomanLorens/logger/log"
)

func main() {
	logger := log.New(log.WithConfig("test.log").WithMaxSize(20).WithMaxAge(7).Build())
	logger2 := log.New(log.WithConfig("test-2.log").Build())

	logger.Info(context.Background(), "boom")
	logger.Debug(context.Background(), "debug %v", 23)
	logger2.Info(context.Background(), "test 2")
}
