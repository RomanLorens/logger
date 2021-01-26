Usage<br/>
logger := log.New(log.WithConfig("test.log").Build())<br />
logger.Info(context.Background(), "boom")<br />
logger.Debug(context.Background(), "debug %v", 23)<br />
