Usage<br/>
logger := log.New(log.WithConfig("test.log").WithMaxSize(20).WithMaxAge(7).Build())<br />
logger.Info(context.Background(), "boom")<br />
logger.Debug(context.Background(), "debug %v", 23)<br />
