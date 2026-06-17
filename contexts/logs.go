package contexts

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logs struct {
	logger *logrus.Logger
}

func NewLogs() *Logs {
	logger := logrus.New()

	logger.Out = os.Stdout

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	logger.SetLevel(logrus.DebugLevel)

	return &Logs{logger: logger}
}

func (l *Logs) Info(v ...any) {
	l.logger.Info(v...)
}

func (l *Logs) Error(v ...any) {
	l.logger.Error(v...)
}

func (l *Logs) Warn(v ...any) {
	l.logger.Warn(v...)
}

func (l *Logs) Debug(msg string, fields ...map[string]any) {
	if len(fields) == 0 || fields[0] == nil || len(fields[0]) == 0 {
		l.logger.Debug(msg)
		return
	}

	l.logger.WithFields(logrus.Fields(fields[0])).Debug(msg)
}

func (l *Logs) InfoFields(msg string, fields map[string]any) {
	l.logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l *Logs) ErrorFields(msg string, fields map[string]any) {
	l.logger.WithFields(logrus.Fields(fields)).Error(msg)
}

func (l *Logs) WarnFields(msg string, fields map[string]any) {
	l.logger.WithFields(logrus.Fields(fields)).Warn(msg)
}
