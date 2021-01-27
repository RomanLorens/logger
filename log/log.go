package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	lj "gopkg.in/natefinch/lumberjack.v2"
)

type contextKey string

//Logger logger
type Logger struct {
	_l *log.Logger
}

const (
	//ReqID request id
	ReqID contextKey = "reqID"
	//UserKey user key
	UserKey contextKey = "user"
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
func New(cfg *Config) *Logger {
	lumber := &lj.Logger{
		Filename:   cfg.LogPath,
		MaxSize:    cfg.MaxSize, // megabytes
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge, // days
	}
	f, err := os.OpenFile(cfg.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Could not create/open log file, %v", err))
	}
	defer f.Close()
	_logger := log.New(f, "", log.Ldate|log.Ltime)
	mw := io.MultiWriter(os.Stdout, lumber)
	_logger.SetOutput(mw)
	out := &Logger{_l: _logger}
	out.Info(context.Background(), "Initialized logger with configuration %v", cfg)
	return out
}

//Info info
func (l Logger) Info(ctx context.Context, format string, args ...interface{}) {
	l._log(ctx, "INFO", format, args...)
}

//Error error
func (l Logger) Error(ctx context.Context, format string, args ...interface{}) {
	l._log(ctx, "ERROR", format, args...)
}

//Warning warning
func (l Logger) Warning(ctx context.Context, format string, args ...interface{}) {
	l._log(ctx, "WARN", format, args...)
}

//Debug debug
func (l Logger) Debug(ctx context.Context, format string, args ...interface{}) {
	l._log(ctx, "DEBUG", format, args...)
}

func (l Logger) _log(ctx context.Context, level string, msg string, args ...interface{}) {
	var user, req string
	user, _ = ctx.Value(UserKey).(string)
	req, _ = ctx.Value(ReqID).(string)
	formatter := "|%v|%v|%v|%v"
	m := fmt.Sprintf(msg, args...)
	l._l.Printf(formatter, strings.ToLower(user), req, level, m)
}

//Panicf panics
func (l Logger) Panicf(ctx context.Context, format string, arg ...interface{}) {
	l.Error(ctx, format, arg...)
	panic(fmt.Sprintf(format, arg...))
}
