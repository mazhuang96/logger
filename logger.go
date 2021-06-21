/**************************************
 * @Author: mazhuang
 * @Date: 2021-06-18 16:32:17
 * @LastEditTime: 2021-06-21 15:31:08
 * @Description:
 **************************************/

package logger

import (
	"fmt"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const defaultTime = "2006/01/02 - 15:04:05.000"

// Logger ...
type Logger struct {
	*zap.Logger

	dir        string
	prefix     string
	format     string
	console    bool
	showline   bool
	stackTrace string
	config     zapcore.EncoderConfig
	level      zapcore.Level
	writer     zapcore.WriteSyncer
}

// NewLogger ...
func NewLogger(level, outDir, prefix string) (log *Logger, err error) {
	log = &Logger{
		dir:     outDir,
		prefix:  prefix,
		console: true,
		format:  "console",
		level:   convertLevel(level),
	}
	if log.prefix != "" {
		log.prefix += " "
	}
	log.initEncoderConfig()

	// 使用file-rotatelogs进行日志分割
	if err = log.getWriteSyncer(); err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	log.Logger = zap.New(log.getEncoderCore(), zap.AddStacktrace(zap.ErrorLevel))
	return
}

// SetShowline configures the Logger to annotate each message with the filename, line number, and function name of zap's caller.
func (l *Logger) SetShowline() {
	l.Logger = l.WithOptions(zap.AddCaller())
}

// SetStacktraceLevel configures the Logger to record a stack trace for all messages at or above a given level.
func (l *Logger) SetStacktraceLevel(level string) {
	l.Logger = l.WithOptions(zap.AddStacktrace(convertLevel(level)))
}

// CloseStacktrace ...
func (l *Logger) CloseStacktrace() {
	l.config.StacktraceKey = ""
	l.wrapCore()
}

// SetTimeFormat sets the log output format.
// default time format is `2006/01/02 - 15:04:05.000`,
func (l *Logger) SetTimeFormat(timeFormat string) {
	if timeFormat == "" {
		return
	}
	l.config.EncodeTime = customTimeEncoder(l.prefix + timeFormat)
	l.wrapCore()
}

// CloseColor ...
func (l *Logger) CloseColor() {
	l.config.EncodeLevel = zapcore.CapitalLevelEncoder
	l.wrapCore()
}

// SetJSONStyle ...
func (l *Logger) SetJSONStyle() {
	l.format = "json"
	l.wrapCore()
}

// wrapCore wraps or replaces the Logger's underlying zapcore.Core.
func (l *Logger) wrapCore() {
	l.Logger = l.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return l.getEncoderCore()
	}))
}

// getWriteSyncer zap logger中加入file-rotatelogs
func (l *Logger) getWriteSyncer() error {
	fileWriter, err := rotatelogs.New(
		path.Join(l.dir, "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		// rotatelogs.WithRotationCount(30),
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

// initEncoderConfig init zapcore.EncoderConfig
func (l *Logger) initEncoderConfig() {
	l.config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     customTimeEncoder(l.prefix + defaultTime),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// getEncoderCore uses the new config to get the new core
func (l *Logger) getEncoderCore() (core zapcore.Core) {
	var encoder zapcore.Encoder
	if l.format == "json" {
		encoder = zapcore.NewJSONEncoder(l.config)
	} else {
		encoder = zapcore.NewConsoleEncoder(l.config)
	}
	return zapcore.NewCore(encoder, l.writer, l.level)
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
