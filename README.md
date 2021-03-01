Usage<br/>
logger, _ := log.New(log.WithConfig("test.log").<br/>
		WithMaxSize(10).WithMaxBackups(8).Build())<br/>
logger2, _ := log.New(log.WithConfig("test-2.log").Build())<br/>

pl, err := log.New(log.WithConfig("/nonesense/aaa.log").Build())<br/>
if err != nil {<br/>
	fmt.Printf("error %v\n", err)<br/>
}<br/>
logger.Info(context.Background(), "boom")<br/>
logger.Debug(context.Background(), "debug %v", 23)<br/>
logger2.Info(context.Background(), "test 2")<br/>
pl.Info(context.Background(), "print logger")<br/>
