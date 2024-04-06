package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var logDirPath = "/tmp/logger_test"

func TestMain(m *testing.M) {
	os.RemoveAll(filepath.Join(logDirPath))
	retCode := m.Run()
	os.Exit(retCode)
}

func TestNewLogger(t *testing.T) {
	_, err := NewLogger(logDirPath, true)
	if err != nil {
		t.Errorf("NewLogger(%s, true) = %v", logDirPath, err)
	}
}

func TestDebug(t *testing.T) {
	logger, _ := NewLogger(logDirPath, true)
	logger.Debug("debug message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "debug.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "log:debug message") {
		t.Errorf("Debug() = %v", messageStr)
	}
}

func TestInfo(t *testing.T) {
	logger, _ := NewLogger(logDirPath, true)
	logger.Info("info message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "info.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "log:info message") {
		t.Errorf("Info() = %v", messageStr)
	}
}

func TestError(t *testing.T) {
	logger, _ := NewLogger(logDirPath, true)
	logger.Error("error message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "error.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "log:error message") {
		t.Errorf("Error() = %v", messageStr)
	}
}

func TestDebugWithFalse(t *testing.T) {
	logger, _ := NewLogger(logDirPath, false)
	logger.Debug("debug message false")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "debug.log"))
	messageStr := string(messages)

	if strings.Contains(messageStr, "log:debug message false") {
		t.Errorf("Debug() = %v", messageStr)
	}
}

func TestSetDebug(t *testing.T) {
	logger, _ := NewLogger(logDirPath, false)
	logger.SetDebug(true)
	logger.Debug("debug message2")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "debug.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "log:debug message2") {
		t.Errorf("Debug() = %v", messageStr)
	}

	logger.SetDebug(false)
	logger.Debug("debug message3")

	messages, _ = os.ReadFile(filepath.Join(logDirPath, "debug.log"))
	messageStr = string(messages)

	if strings.Contains(messageStr, "log:debug message3") {
		t.Errorf("Debug() = %v", messageStr)
	}
}
