package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Цвета для терминала
var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
)

// Уровни логов
const (
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
	LevelDebug = "DEBUG"
)

type Logger struct {
	logger *log.Logger
	level  string
}

func New() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0), // убрал timestamp по дефолту
		level:  LevelInfo,
	}
}

func (l *Logger) log(level, color, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf(msg, args...)
	output := fmt.Sprintf("%s[%s] %s%-5s%s %s",
		Gray, timestamp,
		color, level, Reset,
		formatted,
	)
	l.logger.Println(output)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(LevelInfo, Green, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(LevelWarn, Yellow, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(LevelError, Red, msg, args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(LevelDebug, Cyan, msg, args...)
}
