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
	logger, err := NewLogger(logDirPath, logFileName)
	if err != nil {
		t.Errorf("NewLogger(%s, %s) = %v", logDirPath, logFileName, err)
	}
	if logger == nil {
		t.Errorf("NewLogger(%s, %s) = %v", logDirPath, logFileName, logger)
	}
	if logger.logLevel != "INFO" {
		t.Errorf("NewLogger(%s, %s) = %v", logDirPath, logFileName, logger.logLevel)
	}
}

func TestSetLogLevel(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName)
	logLevels := []string{"DEBUG", "INFO", "ERROR", "NONE"}

	for _, logLevel := range logLevels {
		err := logger.SetLogLevel(logLevel)
		if err != nil {
			t.Errorf("SetLogLevel(%s) = %v", logLevel, err)
		}
		if logger.logLevel != logLevel {
			t.Errorf("SetLogLevel(%s) = %v", logLevel, logger.logLevel)
		}
	}
}

func TestSetLogLevelWithInvalidLogLevel(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName)
	logLevel := "INVALID"
	err := logger.SetLogLevel(logLevel)
	if err == nil {
		t.Errorf("SetLogLevel(%s) = %v", logLevel, err)
	}
	if err.Error() != "ErrInvalidLogLevel" {
		t.Errorf("SetLogLevel(%s) = %v", logLevel, err)
	}
}

func TestDebug(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName)
	logger.SetLogLevel("DEBUG")
	logger.Debug("debug message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, logFileName))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "level:DEBUG\tlog:debug message") {
		t.Errorf("Debug() = %v", messageStr)
	}
}

func TestInfo(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName)
	logger.SetLogLevel("INFO")
	logger.Info("info message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, logFileName))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "level:INFO\tlog:info message") {
		t.Errorf("Info() = %v", messageStr)
	}
}

func TestError(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName)
	logger.SetLogLevel("ERROR")
	logger.Error("error message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, logFileName))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "level:ERROR\tlog:error message") {
		t.Errorf("Error() = %v", messageStr)
	}
}

func TestNone(t *testing.T) {
	logger, _ := NewLogger(logDirPath, logFileName)
	logger.SetLogLevel("NONE")
	logger.Debug("debug message NONE")
	logger.Info("info message NONE")
	logger.Error("error message NONE")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, logFileName))
	messageStr := string(messages)

	if strings.Contains(messageStr, "NONE") {
		t.Errorf("None() = %v", messageStr)
	}
}
