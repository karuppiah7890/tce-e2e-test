package log

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	logDir := createDirectoryIfNotExists()
	writerSync := getLogWriter(logDir)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writerSync, zapcore.DebugLevel)
	globalLogger := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(globalLogger)
}

func createDirectoryIfNotExists() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error while trying to get working directory: %v", err)
	}

	logDir := filepath.Join(path, time.Now().Format("2006-01-02-logs"))

	info, err := os.Stat(logDir)

	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error while trying to get stats about '%s': %v", logDir, err)
		}

		err = os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Error while trying to create logs directory at '%s': %v", logDir, err)
		}
	}

	if info != nil && !info.IsDir() {
		log.Fatalf("Error occurred as we want '%s' to be a directory but it is currently a file", logDir)
	}

	return logDir
}

func getLogWriter(logDirectory string) zapcore.WriteSyncer {
	logFilePath := filepath.Join(logDirectory, "e2e.log")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error while trying to create log file at '%s': %v", logFilePath, err)
	}
	return zapcore.AddSync(logFile)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
	})
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
