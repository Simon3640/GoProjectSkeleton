package providers

import (
	"fmt"
	"log"
	"log/slog"

	contracts "gingoskeleton/src/application/contracts"
)

type LoggerProvider struct {
	enableLog bool
	debugLog  bool
}

func (lp *LoggerProvider) Setup(enableLog bool, debugLog bool) {
	lp.enableLog = enableLog
	lp.debugLog = debugLog
}

var _ contracts.ILoggerProvider = (*LoggerProvider)(nil)

func NewLoggerProvider() *LoggerProvider {
	return &LoggerProvider{}
}

func (lp *LoggerProvider) Error(message string, err error) {
	if lp.enableLog {
		slog.Error(fmt.Sprintf("Error: %s, %s", message, err))
	}
}

func (lp *LoggerProvider) Panic(message string, err error) {
	log.Panicf("Panic: %s, %s", message, err)
}

func (lp *LoggerProvider) ErrorMsg(message string) {
	if lp.enableLog {
		slog.Error(message)
	}
}

func (lp *LoggerProvider) Info(message string) {
	if lp.enableLog {
		slog.Info(message)
	}
}

func (lp *LoggerProvider) Warning(message string) {
	if lp.enableLog {
		slog.Warn(message)
	}
}

func (lp *LoggerProvider) Debug(message string, data any) {
	if lp.enableLog && lp.debugLog {
		if data != nil {
			slog.Debug(fmt.Sprintf("%s: %v", message, data))
		} else {
			slog.Debug(message)
		}
	}
}

var Logger *LoggerProvider

func init() {
	fmt.Println("LoggerProvider initialized")
	Logger = NewLoggerProvider()
}
