package logger

import (
	"context"
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
}

func New() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "[SAYANOX] ", log.LstdFlags),
	}
}

func (l *Logger) Name() string {
	return "Logger"
}

func (l *Logger) Start(ctx context.Context) error {
	l.logger.Println("Logger service started")
	return nil
}

func (l *Logger) Stop(ctx context.Context) error {
	l.logger.Println("Logger service stopped")
	return nil
}

func (l *Logger) Info(msg string) {
	l.logger.Println("[INFO]", msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Println("[ERROR]", msg)
}
