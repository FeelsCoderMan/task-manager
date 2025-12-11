package task

import (
	"log"
	"os"
)

const (
	customLogPrefix = "custom-task-service"
	flags           = log.Ldate | log.Lshortfile
)

type TaskLogger interface {
	error(v any)
	info(v any)
}

type taskLogger struct {
	logger *log.Logger
}

func NewTaskLogger() TaskLogger {
	logger := log.New(os.Stdout, customLogPrefix, flags)
	return taskLogger{
		logger: logger,
	}
}

func (t taskLogger) error(v any) {
	t.logger.SetPrefix(customLogPrefix + " ERROR: ")
	t.logger.Print(v)
}

func (t taskLogger) info(v any) {
	t.logger.SetPrefix(customLogPrefix + " INFO: ")
	t.logger.Print(v)
}
