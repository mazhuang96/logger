/**************************************
 * @Author: mazhuang
 * @Date: 2021-07-01 15:42:21
 * @LastEditTime: 2021-07-06 18:26:49
 * @Description:
 **************************************/

package logger

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Writer ...
type Writer struct {
	logFunc func(msg string, fields ...zapcore.Field)
}

// NewGinLogger ...
func NewGinLogger(log *zap.Logger) *Writer {
	logger := log.WithOptions(zap.AddCallerSkip(3))
	return &Writer{logger.Info}
}

func (l *Writer) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	l.logFunc(string(p))
	return len(p), nil
}

// GinFormatter ...
func GinFormatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("%3d | %13v | %15s | %-7s %#v\n%s",
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}
