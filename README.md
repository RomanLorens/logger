Usage
logger := log.New(log.WithConfig("test.log").Build())
logger.Info(context.Background(), "boom")
logger.Debug(context.Background(), "debug %v", 23)
