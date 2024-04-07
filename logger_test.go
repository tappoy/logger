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
	logger.Debug("message%d", 1)

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "debug.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "debug:message1") {
		t.Errorf("Debug() = %v", messageStr)
	}
}

func TestInfo(t *testing.T) {
	logger, _ := NewLogger(logDirPath, true)
	logger.Info("message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "info.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "info:message") {
		t.Errorf("Info() = %v", messageStr)
	}
}

func TestNotice(t *testing.T) {
	logger, _ := NewLogger(logDirPath, true)
	logger.Notice("message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "notice.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "notice:message") {
		t.Errorf("Notice() = %v", messageStr)
	}
}

func TestError(t *testing.T) {
	logger, _ := NewLogger(logDirPath, true)
	logger.Error("message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "error.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "error:message") {
		t.Errorf("Error() = %v", messageStr)
	}
}

func TestDebugWithFalse(t *testing.T) {
	logger, _ := NewLogger(logDirPath, false)
	logger.Debug("debug message false")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "debug.log"))
	messageStr := string(messages)

	if strings.Contains(messageStr, "debug:debug message false") {
		t.Errorf("Debug() = %v", messageStr)
	}
}

func TestSetDebug(t *testing.T) {
	logger, _ := NewLogger(logDirPath, false)
	logger.SetDebug(true)
	logger.Debug("debug message2")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "debug.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "debug:debug message2") {
		t.Errorf("Debug() = %v", messageStr)
	}

	logger.SetDebug(false)
	logger.Debug("debug message3")

	messages, _ = os.ReadFile(filepath.Join(logDirPath, "debug.log"))
	messageStr = string(messages)

	if strings.Contains(messageStr, "debug:debug message3") {
		t.Errorf("Debug() = %v", messageStr)
	}
}
