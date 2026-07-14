package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/hamidgh01/cache-proxy/config"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func parseLevel(level string) Level {
	switch level {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warning":
		return WARNING
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

type Logger struct {
	*log.Logger
	level Level
}

func NewLogger(cfg config.LoggerConf) *Logger {

	output := os.Stdout
	if cfg.OutputFile != "" {
		file, err := os.OpenFile(cfg.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			output = file
		}
	}

	return &Logger{
		Logger: log.New(output, "", log.LstdFlags),
		level:  parseLevel(cfg.Level),
	}
}

func (l *Logger) Debug(message string) {
	if l.level <= DEBUG {
		l.Printf("[DEBUG] %s\n", message)
	}
}

func (l *Logger) Debugf(message string, v ...any) {
	if l.level <= DEBUG {
		l.Printf("[DEBUG] %s\n", fmt.Sprintf(message, v...))
	}
}

func (l *Logger) Info(message string) {
	if l.level <= INFO {
		l.Printf("[INFO] %s\n", message)
	}
}

func (l *Logger) Infof(message string, v ...any) {
	if l.level <= INFO {
		l.Printf("[INFO] %s\n", fmt.Sprintf(message, v...))
	}
}

func (l *Logger) Warning(message string) {
	if l.level <= WARNING {
		l.Printf("[WARNING] %s\n", message)
	}
}

func (l *Logger) Warningf(message string, v ...any) {
	if l.level <= WARNING {
		l.Printf("[WARNING] %s\n", fmt.Sprintf(message, v...))
	}
}

func (l *Logger) Error(message string) {
	if l.level <= ERROR {
		l.Printf("[ERROR] %s\n", message)
	}
}

func (l *Logger) Errorf(message string, v ...any) {
	if l.level <= ERROR {
		l.Printf("[ERROR] %s\n", fmt.Sprintf(message, v...))
	}
}

func (l *Logger) Fatal(message string) {
	if l.level <= FATAL {
		l.Printf("[FATAL] %s\n", message)
		os.Exit(1)
	}
}

func (l *Logger) Fatalf(message string, v ...any) {
	if l.level <= FATAL {
		l.Printf("[FATAL] %s\n", fmt.Sprintf(message, v...))
		os.Exit(1)
	}
}
