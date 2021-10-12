/**************************************
 * @Author: mazhuang
 * @Date: 2021-06-18 16:32:17
 * @LastEditTime: 2021-10-12 10:44:50
 * @Description:
 **************************************/

package logger

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const defaultTime = "2006/01/02 - 15:04:05.000"

// Logger ...
type Logger struct {
	*zap.Logger

	format  string // encoder format
	console bool
	config  zapcore.EncoderConfig
	level   zapcore.Level
	writer  zapcore.WriteSyncer
}

// Config ...
type Config struct {
	Level      string // the lowest level of log printing , default "info"
	Dir        string // log output folder
	Prefix     string // prefix on each line to identify the logger
	TimeFormat string // the time format of each line in the log
	MaxAge     int    // max age (days) of each log file, default 7 days
	Color      bool   // the color of log level
	ShowLine   bool   // show log call line number
	Stacktrace string // stack trace log level
	Encoder    string // log encoding format, divided into "json" and "console", default "console"
}

// New ...
func New(c Config) (log *Logger, err error) {
	log = &Logger{
		format:  c.Encoder,
		console: true,
		level:   convertLevel(c.Level),
	}
	if c.Prefix != "" {
		c.Prefix += " "
	}
	if c.TimeFormat == "" {
		c.TimeFormat = defaultTime
	}
	log.config = initEncoderConfig(c.Prefix + c.TimeFormat)
	if c.Stacktrace == "" {
		log.config.StacktraceKey = ""
	}
	if !c.Color {
		// close color
		log.config.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	if err = log.setWriter(c.Dir, c.MaxAge); err != nil {
		fmt.Printf("set writer failed: %v", err.Error())
		return
	}

	log.Logger = zap.New(log.getEncoderCore(log.config), zap.AddStacktrace(convertLevel(c.Stacktrace)))
	if c.ShowLine {
		log.Logger = log.WithOptions(zap.AddCaller())
	}
	return
}

// NewDefault ...
func NewDefault() (log *Logger, err error) {
	return New(Config{
		Color:      true,
		Dir:        "logs",
		MaxAge:     7,
		TimeFormat: defaultTime,
		Stacktrace: "error",
		ShowLine:   false,
		Encoder:    "console",
	})
}

// GinLogConfig return gin web framework log configuration
func (l *Logger) GinLogConfig() gin.LoggerConfig {
	return gin.LoggerConfig{
		Output:    NewGinLogger(l.Logger),
		Formatter: GinFormatter,
	}
}

// Showline configures the Logger to annotate each message with the filename, line number, and function name of zap's caller.
func (l *Logger) Showline() *zap.Logger {
	return l.WithOptions(zap.AddCaller())
}

// SetStacktraceLevel configures the Logger to record a stack trace for all messages at or above a given level.
func (l *Logger) SetStacktraceLevel(level string) *zap.Logger {
	return l.WithOptions(zap.AddStacktrace(convertLevel(level)))
}

// CloseStacktrace ...
func (l *Logger) CloseStacktrace() *zap.Logger {
	c := l.config
	c.StacktraceKey = ""
	return l.wrapCore(c)
}

// SetTimeFormat sets the log output format.
// default time format is `2006/01/02 - 15:04:05.000`,
func (l *Logger) SetTimeFormat(timeFormat string) *zap.Logger {
	c := l.config
	c.EncodeTime = customTimeEncoder(timeFormat)
	return l.wrapCore(c)
}

// NoColor ...
func (l *Logger) NoColor() *zap.Logger {
	c := l.config
	c.EncodeLevel = zapcore.CapitalLevelEncoder
	return l.wrapCore(c)
}

// SetJSONStyle change output style to json
func (l *Logger) SetJSONStyle() *zap.Logger {
	tmp := l.format
	defer func() { l.format = tmp }()
	l.format = "json"
	return l.wrapCore(l.config)
}

// wrapCore wraps or replaces the Logger's underlying zapcore.Core.
func (l *Logger) wrapCore(ec zapcore.EncoderConfig) *zap.Logger {
	return l.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return l.getEncoderCore(ec)
	}))
}

// setWriter zap logger's writer use file-rotatelogs
func (l *Logger) setWriter(dir string, maxAge int) error {
	if maxAge <= 0 {
		maxAge = 7
	}
	fileWriter, err := rotatelogs.New(
		path.Join(dir, "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(time.Duration(maxAge)*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return err
	}
	if l.console {
		l.writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
		return nil
	}
	l.writer = zapcore.AddSync(fileWriter)
	return nil
}

// getEncoderCore uses the new config to get the new core
func (l *Logger) getEncoderCore(ec zapcore.EncoderConfig) (core zapcore.Core) {
	var encoder zapcore.Encoder
	if l.format == "json" {
		encoder = zapcore.NewJSONEncoder(ec)
	} else {
		encoder = zapcore.NewConsoleEncoder(ec)
	}
	return zapcore.NewCore(encoder, l.writer, l.level)
}

// initEncoderConfig init zapcore.EncoderConfig
func initEncoderConfig(format string) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     customTimeEncoder(format),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// customTimeEncoder sets custom log output time format
func customTimeEncoder(format string) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(format))
	}
}

func convertLevel(lvl string) zapcore.Level {
	switch lvl {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
