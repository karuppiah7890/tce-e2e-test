package log

import "go.uber.org/zap"

func Fatalf(template string, args ...interface{}) {
	zap.S().Fatalf(template, args...)
}

func Fatal(args ...interface{}) {
	zap.S().Fatal(args...)
}

func Infof(template string, args ...interface{}) {
	zap.S().Infof(template, args...)
}

func Info(args ...interface{}) {
	zap.S().Info(args...)
}

func Warnf(template string, args ...interface{}) {
	zap.S().Warnf(template, args...)
}

func Warn(args ...interface{}) {
	zap.S().Warn(args...)
}

func DPanicf(template string, args ...interface{}) {
	zap.S().DPanicf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	zap.S().Errorf(template, args...)
}

// infoWriter implements the io.Writer interface
type infoWriter struct{}

func (s infoWriter) Write(p []byte) (n int, err error) {
	zap.S().Info(string(p))
	return len(p), nil
}

var InfoWriter = infoWriter{}

// errorWriter implements the io.Writer interface
type errorWriter struct{}

func (s errorWriter) Write(p []byte) (n int, err error) {
	zap.S().Error(string(p))
	return len(p), nil
}

var ErrorWriter = errorWriter{}
