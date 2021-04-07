package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	lj "gopkg.in/natefinch/lumberjack.v2"
)

type contextKey string

type fileLogger struct {
	_l          *log.Logger
	withLogName bool
	format      string
}

//Logger logger
type Logger interface {
	Info(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, format string, args ...interface{})
	Warning(ctx context.Context, format string, args ...interface{})
	Debug(ctx context.Context, format string, args ...interface{})
	Panicf(ctx context.Context, format string, arg ...interface{})
}

const (
	//ReqID request id
	ReqID contextKey = "reqID"
	//UserKey user key
	UserKey contextKey = "user"
	//LogName logger name
	LogName contextKey = "logName"
)

var (
	fileFormatter  = "|%v|%v|%v|%v"
	printFormatter = "%v |%v|%v|%v|%v"
)

//Config configuration
type Config struct {
	LogPath string
	//MaxSize in megabytes
	MaxSize int
	//MaxBackups
	MaxBackups int
	//MaxAge in days
	MaxAge int
	//LogName log name optional
	LogName bool
}

//ConfigBuilder builder
type ConfigBuilder struct {
	Config
}

type kv struct {
	Key   string
	Value int
}

//WithConfig config builder
func WithConfig(logPath string) *ConfigBuilder {
	return &ConfigBuilder{
		Config: Config{LogPath: logPath},
	}
}

//Build builds config
func (b ConfigBuilder) Build() *Config {
	c := &Config{LogPath: b.LogPath, MaxAge: b.MaxAge, MaxBackups: b.MaxBackups, MaxSize: b.MaxSize, LogName: b.LogName}
	if c.MaxAge == 0 {
		c.MaxAge = 7
	}
	if c.MaxSize == 0 {
		c.MaxSize = 10
	}
	if c.MaxBackups == 0 {
		c.MaxBackups = 3
	}
	return c
}

//WithMaxSize max size in MB
func (b *ConfigBuilder) WithMaxSize(s int) *ConfigBuilder {
	b.MaxSize = s
	return b
}

//WithLogName format tokens
func (b *ConfigBuilder) WithLogName(l bool) *ConfigBuilder {
	b.LogName = l
	return b
}

//WithMaxAge max age
func (b *ConfigBuilder) WithMaxAge(s int) *ConfigBuilder {
	b.MaxAge = s
	return b
}

//WithMaxBackups max backups
func (b *ConfigBuilder) WithMaxBackups(s int) *ConfigBuilder {
	b.MaxBackups = s
	return b
}

//New inits logger
func New(cfg *Config) (Logger, error) {
	lumber := &lj.Logger{
		Filename:   cfg.LogPath,
		MaxSize:    cfg.MaxSize, // megabytes
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge, // days
	}
	f, err := os.OpenFile(cfg.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Could not create/open log file - will use print logger, %v\n", err)
		return PrintLogger(cfg.LogName), err
	}
	defer f.Close()
	_logger := log.New(f, "", log.Ldate|log.Ltime)
	mw := io.MultiWriter(os.Stdout, lumber)
	_logger.SetOutput(mw)

	format := fileFormatter
	if cfg.LogName {
		format += "|%v"
	}
	out := &fileLogger{_l: _logger, withLogName: cfg.LogName, format: format}
	out.Info(context.Background(), "Initialized logger with configuration %v", cfg)
	logUncaughtPanic(f, out)
	return out, nil
}

//Info info
func (l fileLogger) Info(ctx context.Context, format string, args ...interface{}) {
	_log(ctx, l, "INFO", format, args...)
}

//Error error
func (l fileLogger) Error(ctx context.Context, format string, args ...interface{}) {
	_log(ctx, l, "ERROR", format, args...)
}

//Warning warning
func (l fileLogger) Warning(ctx context.Context, format string, args ...interface{}) {
	_log(ctx, l, "WARN", format, args...)
}

//Debug debug
func (l fileLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	_log(ctx, l, "DEBUG", format, args...)
}

//Panicf panics
func (l fileLogger) Panicf(ctx context.Context, format string, arg ...interface{}) {
	l.Error(ctx, format, arg...)
	panic(fmt.Sprintf(format, arg...))
}

func _log(ctx context.Context, l Logger, level string, msg string, args ...interface{}) {
	m := fmt.Sprintf(msg, args...)
	var user, req, logName string
	user, _ = ctx.Value(UserKey).(string)
	req, _ = ctx.Value(ReqID).(string)
	logName, _ = ctx.Value(LogName).(string)
	switch v := l.(type) {
	case fileLogger:
		if v.withLogName {
			v._l.Printf(v.format, strings.ToLower(user), req, level, logName, m)
		} else {
			v._l.Printf(v.format, strings.ToLower(user), req, level, m)
		}
	case printLogger:
		date := time.Now().Format("2006/01/01 15:04:05")
		if v.withLogName {
			fmt.Printf(v.format+"\n", date, strings.ToLower(user), req, level, logName, m)
		} else {
			fmt.Printf(v.format+"\n", date, strings.ToLower(user), req, level, m)
		}
	}
}
