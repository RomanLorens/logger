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
	_l *log.Logger
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
	UserKey        contextKey = "user"
	fileFormatter             = "|%v|%v|%v|%v"
	printFormatter            = "%v |%v|%v|%v|%v"
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
}

//ConfigBuilder builder
type ConfigBuilder struct {
	Config
}

//WithConfig config builder
func WithConfig(logPath string) *ConfigBuilder {
	return &ConfigBuilder{
		Config: Config{LogPath: logPath},
	}
}

//Build builds config
func (b ConfigBuilder) Build() *Config {
	c := &Config{LogPath: b.LogPath, MaxAge: b.MaxAge, MaxBackups: b.MaxBackups, MaxSize: b.MaxSize}
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
		return &printLogger{}, err
	}
	defer f.Close()
	_logger := log.New(f, "", log.Ldate|log.Ltime)
	mw := io.MultiWriter(os.Stdout, lumber)
	_logger.SetOutput(mw)
	out := &fileLogger{_l: _logger}
	out.Info(context.Background(), "Initialized logger with configuration %v", cfg)
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
	var user, req string
	user, _ = ctx.Value(UserKey).(string)
	req, _ = ctx.Value(ReqID).(string)
	m := fmt.Sprintf(msg, args...)
	switch v := l.(type) {
	case fileLogger:
		v._l.Printf(fileFormatter, strings.ToLower(user), req, level, m)
	case printLogger:
		date := time.Now().Format("2006/01/01 15:04:05")
		fmt.Printf(printFormatter+"\n", date, strings.ToLower(user), req, level, m)
	}

}
