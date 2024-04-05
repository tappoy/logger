package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var logDirPath = "/tmp"
var logFileName = "test.log"

func TestMain(m *testing.M) {
	os.RemoveAll(filepath.Join(logDirPath, logFileName))
	retCode := m.Run()
	os.Exit(retCode)
}

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger(logDirPath, logFileName, "INFO")
	if err != nil {
		t.Errorf("NewLogger(%s, %s) = %v", logDirPath, logFileName, err)
	}
	if logger == nil {
		t.Errorf("NewLogger(%s, %s) = %v", logDirPath, logFileName, logger)
	}
	if logger.level != "INFO" {
		t.Errorf("NewLogger(%s, %s) = %v", logDirPath, logFileName, logger.level)
	}
}

func TestNewLoggerWithInvalidLogLevel(t *testing.T) {
  _, err := NewLogger(logDirPath, logFileName, "INVALID")
  if err == nil {
    t.Errorf("NewLogger(%s, %s) = %v", logDirPath, logFileName, err)
  }
  if err.Error() != "ErrInvalidLogLevel" {
    t.Errorf("NewLogger(%s, %s) = %v", logDirPath, logFileName, err)
  }
}

func TestSetLogLevel(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName, "INFO")
	logLevels := []string{"DEBUG", "INFO", "ERROR", "NONE"}

	for _, level := range logLevels {
		err := logger.SetLogLevel(level)
		if err != nil {
			t.Errorf("SetLogLevel(%s) = %v", level, err)
		}
		if logger.level != level {
			t.Errorf("SetLogLevel(%s) = %v", level, logger.level)
		}
	}
}

func TestSetLogLevelWithInvalidLogLevel(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName, "INFO")
	level := "INVALID"
	err := logger.SetLogLevel(level)
	if err == nil {
		t.Errorf("SetLogLevel(%s) = %v", level, err)
	}
	if err.Error() != "ErrInvalidLogLevel" {
		t.Errorf("SetLogLevel(%s) = %v", level, err)
	}
}

func TestDebug(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName, "DEBUG")
	logger.Debug("debug message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, logFileName))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "level:DEBUG\tlog:debug message") {
		t.Errorf("Debug() = %v", messageStr)
	}
}

func TestInfo(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName, "INFO")
	logger.Info("info message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, logFileName))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "level:INFO\tlog:info message") {
		t.Errorf("Info() = %v", messageStr)
	}
}

func TestError(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName, "ERROR")
	logger.Error("error message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, logFileName))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "level:ERROR\tlog:error message") {
		t.Errorf("Error() = %v", messageStr)
	}
}

func TestNone(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName, "NONE")
	logger.Debug("debug message NONE")
	logger.Info("info message NONE")
	logger.Error("error message NONE")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, logFileName))
	messageStr := string(messages)

	if strings.Contains(messageStr, "NONE") {
		t.Errorf("None() = %v", messageStr)
	}
}
