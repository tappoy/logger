package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var logDir = "/tmp/logger_test"

func TestMain(m *testing.M) {
	os.RemoveAll(filepath.Join(logDir))
	retCode := m.Run()
	os.Exit(retCode)
}

func TestNewLogger(t *testing.T) {
	_, err := NewLogger(logDir)
	if err != nil {
		t.Errorf("NewLogger(%s, true) = %v", logDir, err)
	}
}

func TestDebug(t *testing.T) {
	logger, _ := NewLogger(logDir)

	// It should not create and write to debug log file
	{
		logger.Debug("message%d", 1)
		_, err := os.ReadFile(filepath.Join(logDir, "debug.log"))
		if err == nil {
			t.Errorf("Debug() = %v", err)
		}
	}

	// It should create and write to debug log file
	{
		// create debug log file
		err := createFileIfNotExist(logDir, "debug.log")
		if err != nil {
			t.Errorf("Debug() = %v", err)
		}
		// log error
		logger.Debug("message%d", 2)
		messages, _ := os.ReadFile(filepath.Join(logDir, "debug.log"))
		messageStr := string(messages)

		if !strings.Contains(messageStr, "debug:message2") {
			t.Errorf("Debug() = %v", messageStr)
		}
	}

	// remove debug log file
	os.Remove(filepath.Join(logDir, "debug.log"))

	// It should not create and write to debug log file
	{
		logger.Debug("message%d", 3)
		_, err := os.ReadFile(filepath.Join(logDir, "debug.log"))
		if err == nil {
			t.Errorf("Debug() = %v", err)
		}
	}

}

func TestInfo(t *testing.T) {
	logger, _ := NewLogger(logDir)
	logger.Info("message")

	messages, _ := os.ReadFile(filepath.Join(logDir, "info.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "info:message") {
		t.Errorf("Info() = %v", messageStr)
	}
}

func TestNotice(t *testing.T) {
	logger, _ := NewLogger(logDir)
	logger.Notice("message")

	messages, _ := os.ReadFile(filepath.Join(logDir, "notice.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "notice:message") {
		t.Errorf("Notice() = %v", messageStr)
	}
}

func TestError(t *testing.T) {
	logger, _ := NewLogger(logDir)
	logger.Error("message")

	messages, _ := os.ReadFile(filepath.Join(logDir, "error.log"))
	messageStr := string(messages)

	if !strings.Contains(messageStr, "error:message") {
		t.Errorf("Error() = %v", messageStr)
	}
}

func TestRotate(t *testing.T) {
	logger, _ := NewLogger(logDir)

	// make test times
	today := time.Now()
	tommorow := today.AddDate(0, 0, 1)

	// add today's last message to log file
	logger.log("info", today, "last message")

	// add tommorow's message to log file
	logger.log("info", tommorow, "rotated")

	// check if log file is rotated
	messages, err := os.ReadFile(filepath.Join(logDir, "info_"+today.Format("2006-01-02")+".log"))
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}

	messageStr := string(messages)
	if !strings.Contains(messageStr, "info:last message") {
		t.Errorf("Rotate() = %v", messageStr)
	}

	// check if new log file is created
	messages, err = os.ReadFile(filepath.Join(logDir, "info.log"))
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}

	messageStr = string(messages)
	if !strings.Contains(messageStr, "info:rotated") {
		t.Errorf("Rotate() = %v", messageStr)
	}

}

func TestRotateFailed(t *testing.T) {
	logger, _ := NewLogger(logDir)

	// make test times
	today := time.Now()
	tommorow := today.AddDate(0, 0, 1)

	// create roteated log file for it's to be failed
	err := createFileIfNotExist(logDir, "rotate-railed_"+today.Format("2006-01-02")+".log")
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}

	// set permission to read only
	err = os.Chmod(filepath.Join(logDir, "rotate-railed_"+today.Format("2006-01-02")+".log"), 0400)
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}

	// add today's last message to log file
	logger.log("rotate-railed", today, "Roteate Failed1")

	// add tommorow's message to log file
	logger.log("rotate-railed", tommorow, "Roteate Failed2")

	// check if log file is not rotated
	messages, err := os.ReadFile(filepath.Join(logDir, "rotate-railed_"+today.Format("2006-01-02")+".log"))
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}

	messageStr := string(messages)
	if strings.Contains(messageStr, "rotate-railed:last message") {
		t.Errorf("Rotate() = %v", messageStr)
	}

	// check both messages are in the info.log
	messages, err = os.ReadFile(filepath.Join(logDir, "rotate-railed.log"))
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}

	messageStr = string(messages)
	if !strings.Contains(messageStr, "rotate-railed:Roteate Failed1") {
		t.Errorf("Rotate() = %v", messageStr)
	}
	if !strings.Contains(messageStr, "rotate-railed:Roteate Failed2") {
		t.Errorf("Rotate() = %v", messageStr)
	}

}
