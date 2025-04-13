package logger

import (
	"fmt"
	"github.com/ttrtcixy/users/internal/logger/lib"
	"log/slog"
	"os"
)

//go:generate mockery --name=Logger

type Logger interface {
	Info(info string)
	Error(err string)
	ErrorOp(err string, op string)
	Fatal(err string)
	Debug(info string)
}

type slogLogger struct {
	log *slog.Logger
}

func Load() Logger {
	opts := colorLog.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := colorLog.NewPrettyHandler(os.Stderr, opts)
	logger := slog.New(handler)

	return slogLogger{log: logger}
}

func (s slogLogger) Info(info string) {
	s.log.Info(info)
}

func (s slogLogger) Error(err string) {
	s.log.Error(err)
}

func (s slogLogger) ErrorOp(err string, op string) {
	s.log.Error(fmt.Sprintf("%s: %s", op, err))
}

func (s slogLogger) Fatal(err string) {
	s.log.Error(err)
	os.Exit(1)
}

func (s slogLogger) Debug(info string) {
	s.log.Debug(info)
}
