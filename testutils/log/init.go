package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: Check what happens when multiple tests run in parallel and call this.
// Does it support multiple parallel tests logging to different files? Or is there any problems
// due to the usage of global loggers and zap.ReplaceGlobals 😅 in which case we have to resort to
// using local loggers along with passing around loggers for usage in functions / methods etc
func InitLogger(loggingProgram string) {
	logDir := createDirectoryIfNotExists()
	writerSync := getLogWriter(logDir, loggingProgram)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writerSync, zapcore.DebugLevel)
	// TODO: The caller is always log/log.go and it's not useful as we don't know which function in the stack called it.
	// Can we stack information etc? Or we will remove it for now
	// Added AddCallerSkip to log stack for above todo
	globalLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

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

func getLogWriter(logDirectory string, loggingProgram string) zapcore.WriteSyncer {
	logFilePath := filepath.Join(logDirectory, fmt.Sprintf("%s-%s.log", loggingProgram, time.Now().Format("2006-01-02-15-04-05")))
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error while trying to create log file at '%s': %v", logFilePath, err)
	}
	// This way we log to standard output and to the log file, kind of like tee command in Linux :D
	stdOutAndLogFile := io.MultiWriter(os.Stdout, logFile)
	return zapcore.AddSync(stdOutAndLogFile)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
	})
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
