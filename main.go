package main

import (
	"context"
	"fmt"

	"github.com/RomanLorens/logger/log"
)

func main() {
	logger, _ := log.New(log.WithConfig("test.log").
		WithMaxSize(10).WithMaxBackups(8).Build())
	logger2, _ := log.New(log.WithConfig("test-2.log").Build())

	pl, err := log.New(log.WithConfig("/nonesense/aaa.log").Build())
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	logger.Info(context.Background(), "boom")
	logger.Debug(context.Background(), "debug %v", 23)
	logger2.Info(context.Background(), "test 2")
	pl.Info(context.Background(), "print logger")
}
