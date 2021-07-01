/**************************************
 * @Author: mazhuang
 * @Date: 2021-07-01 15:42:21
 * @LastEditTime: 2021-07-01 15:50:05
 * @Description:
 **************************************/

package logger

import (
	"bytes"

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
