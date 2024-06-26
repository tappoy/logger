package logger

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

var logDir = "tmp/logger_test"

func TestMain(m *testing.M) {
	syscall.Umask(0)
	os.RemoveAll(filepath.Join(logDir))
	retCode := m.Run()
	os.Exit(retCode)
}

func TestNewLogger(t *testing.T) {
	_, err := NewLogger(logDir)
	if err != nil {
		t.Errorf("NewLogger(%s) = %v", logDir, err)
	}

	// check if log directory is created
	s, err := os.Stat(logDir)
	if err != nil {
		t.Errorf("NewLogger(%s) = %v", logDir, err)
	}
	if !s.IsDir() {
		t.Errorf("NewLogger(%s) = %v", logDir, s)
	}
	// check permission
	if s.Mode().Perm() != 0775 {
		t.Errorf("NewLogger(%s) permission must be 0775. got %v", logDir, s.Mode().Perm())
	}

	// check if temp file is removed
	// get tempfile from log directory by glob
	tempfiles, _ := filepath.Glob(filepath.Join(logDir, "tempfile*"))
	if len(tempfiles) != 0 {
		t.Errorf("NewLogger(%s): tempfiles must be removed. got %v", logDir, tempfiles)
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
		err := createFileIfNotExist(filepath.Join(logDir, "debug.log"))
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

	// rotate log directory
	rotateLogDir := filepath.Join(logDir, "rotate")

	// make test times
	today := time.Now()
	tommorow := today.AddDate(0, 0, 1)

	// add today's last message to log file
	logger.log("info", today, "last message")

	// add tommorow's message to log file
	logger.log("info", tommorow, "rotated")

	// check if log file is rotated
	messages, err := os.ReadFile(filepath.Join(rotateLogDir, today.Format("2006-01-02")+".info.log"))
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

	// rotate log directory
	rotateLogDir := filepath.Join(logDir, "rotate")

	// make test times
	today := time.Now()
	tommorow := today.AddDate(0, 0, 1)

	// rotate file name
	rotateFile := filepath.Join(rotateLogDir, today.Format("2006-01-02")+".rotate-failed.log")

	// create roteated log file for it's to be failed
	err := createFileIfNotExist(rotateFile)
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}

	// set permission to read only
	err = os.Chmod(rotateFile, 0400)
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}
	defer os.Chmod(rotateLogDir, 0775)

	// add today's last message to log file
	logger.log("rotate-failed", today, "Roteate Failed1")

	// add tommorow's message to log file
	logger.log("rotate-failed", tommorow, "Roteate Failed2")

	// check if log file is not rotated
	messages, err := os.ReadFile(rotateFile)
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}

	messageStr := string(messages)
	if strings.Contains(messageStr, "rotate-failed:last message") {
		t.Errorf("Rotate() = %v", messageStr)
	}

	// check both messages are in the info.log
	messages, err = os.ReadFile(filepath.Join(logDir, "rotate-failed.log"))
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}

	messageStr = string(messages)
	if !strings.Contains(messageStr, "rotate-failed:Roteate Failed1") {
		t.Errorf("Rotate() = %v", messageStr)
	}
	if !strings.Contains(messageStr, "rotate-failed:Roteate Failed2") {
		t.Errorf("Rotate() = %v", messageStr)
	}

}

func TestRotateOver30Files(t *testing.T) {
	logger, _ := NewLogger(logDir)

	// rotate log directory
	rotateLogDir := filepath.Join(logDir, "rotate")

	// make test times
	today := time.Now()

	// add 40 log files
	for i := 0; i < 40; i++ {
		psudoDate := today.AddDate(0, 0, i)
		logger.log("info", psudoDate, "message")
		os.Chtimes(filepath.Join(logDir, "info.log"), psudoDate, psudoDate)
	}

	// check rotate directory has 30 files
	files, err := os.ReadDir(rotateLogDir)
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}
	if len(files) != 30 {
		t.Errorf("Rotate() = %v", files)
	}

	// chmod to read only to rotate directory
	err = os.Chmod(rotateLogDir, 0400)
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}
	defer os.Chmod(rotateLogDir, 0775)

	// check buckup directory's mode
	s, err := os.Stat(rotateLogDir)
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	}
	if s.Mode().Perm() != 0400 {
		t.Errorf("Rotate() = %v", s.Mode().Perm())
	}

	testMessage := "This is written to info.log even if a rotate error occurs."

	// add 10 more log files
	for i := 40; i < 50; i++ {
		psudoDate := today.AddDate(0, 0, i)
		logger.log("info", psudoDate, testMessage)
		os.Chtimes(filepath.Join(logDir, "info.log"), psudoDate, psudoDate)
	}

	// check logger.log is created
	_, err = os.ReadFile(filepath.Join(logDir, "logger.log"))
	if err != nil {
		t.Errorf("Rotate() = %v", err)
	} else {
		// check if test message is written to logger.log
		messages, _ := os.ReadFile(filepath.Join(logDir, "info.log"))
		messageStr := string(messages)
		if !strings.Contains(messageStr, testMessage) {
			t.Errorf("Rotate() = %v", messageStr)
		}
	}
}

// When no write permission to log directory
func TestNoWritePermission(t *testing.T) {
	// set permission to read only
	err := os.Chmod(logDir, 0400)
	defer os.Chmod(logDir, 0775)
	if err != nil {
		t.Fatal("Cant set permission to read only")
	}

	_, err = NewLogger(logDir)
	if err == nil {
		t.Errorf("NoWritePermission() = %v", err)
	}
}
