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
	_, err := NewLogger(logDirPath)
	if err != nil {
		t.Errorf("NewLogger(%s, true) = %v", logDirPath, err)
	}
}

func TestDebug(t *testing.T) {
	logger, _ := NewLogger(logDirPath)

	// It should not create and write to debug log file
	{
		logger.Debug("message%d", 1)
		_, err := os.ReadFile(filepath.Join(logDirPath, "debug.log"))
		if err == nil {
			t.Errorf("Debug() = %v", err)
		}
	}

	// It should create and write to debug log file
	{
		// create debug log file
		err := createFileIfNotExist(logDirPath, "debug.log")
		if err != nil {
			t.Errorf("Debug() = %v", err)
		}
		// log error
		logger.Debug("message%d", 2)
		messages, _ := os.ReadFile(filepath.Join(logDirPath, "debug.log"))
		messageStr := string(messages)

		if !strings.Contains(messageStr, "debug:message2") {
			t.Errorf("Debug() = %v", messageStr)
		}
	}

	// remove debug log file
	os.Remove(filepath.Join(logDirPath, "debug.log"))

	// It should not create and write to debug log file
	{
		logger.Debug("message%d", 3)
		_, err := os.ReadFile(filepath.Join(logDirPath, "debug.log"))
		if err == nil {
			t.Errorf("Debug() = %v", err)
		}
	}

}

func TestInfo(t *testing.T) {
	logger, _ := NewLogger(logDirPath)
	logger.Info("message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "info.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "info:message") {
		t.Errorf("Info() = %v", messageStr)
	}
}

func TestNotice(t *testing.T) {
	logger, _ := NewLogger(logDirPath)
	logger.Notice("message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "notice.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "notice:message") {
		t.Errorf("Notice() = %v", messageStr)
	}
}

func TestError(t *testing.T) {
	logger, _ := NewLogger(logDirPath)
	logger.Error("message")

	messages, _ := os.ReadFile(filepath.Join(logDirPath, "error.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "error:message") {
		t.Errorf("Error() = %v", messageStr)
	}
}
