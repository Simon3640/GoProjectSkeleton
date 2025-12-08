package providers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
)

type LoggerProvider struct {
	enableLog bool
	debugLog  bool
	logger    *slog.Logger
}

type LogHandler struct {
	slog.Handler
}

func (h LogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h LogHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h LogHandler) Handle(ctx context.Context, record slog.Record) error {
	var color string
	switch record.Level {
	case slog.LevelDebug:
		color = "\033[36m"
	case slog.LevelInfo:
		color = "\033[32m"
	case slog.LevelWarn:
		color = "\033[33m"
	case slog.LevelError:
		color = "\033[31m"
	default:
		color = "\033[0m"
	}

	if _, err := os.Stdout.WriteString(color + record.Message + "\033[0m" + "\n"); err != nil {
		return err
	}
	return nil
}

func (lp *LoggerProvider) Setup(enableLog bool, debugLog bool) {
	lp.enableLog = enableLog
	lp.debugLog = debugLog
	lp.logger = slog.New(LogHandler{})
}

var _ contractsProviders.ILoggerProvider = (*LoggerProvider)(nil)

func NewLoggerProvider() *LoggerProvider {
	return &LoggerProvider{}
}

func (lp *LoggerProvider) Error(message string, err error) {
	if lp.enableLog {
		lp.logger.Error(fmt.Sprintf("Error: %s, %s", message, err))
	}
}

func (lp *LoggerProvider) Panic(message string, err error) {
	if lp.enableLog {
		log.Panicf("%s", fmt.Sprintf("Panic: %s, %s", message, err))
	}
}

func (lp *LoggerProvider) ErrorMsg(message string) {
	if lp.enableLog {
		lp.logger.Error(message)
	}
}

func (lp *LoggerProvider) Info(message string) {
	if lp.enableLog {
		lp.logger.Info(message)
	}
}

func (lp *LoggerProvider) Warning(message string) {
	if lp.enableLog {
		lp.logger.Warn(message)
	}
}

func (lp *LoggerProvider) Debug(message string, data any) {
	if lp.enableLog && lp.debugLog {
		if data != nil {
			lp.logger.Debug(fmt.Sprintf("%s: %v", message, data))
		} else {
			lp.logger.Debug(message)
		}
	}
}

var Logger *LoggerProvider

func init() {
	fmt.Println("LoggerProvider initialized")
	Logger = NewLoggerProvider()
}
