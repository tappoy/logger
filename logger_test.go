package logger

import (
  "testing"
)

func TestNewLogger(t *testing.T) {
  logDirPath := "/tmp"
  logFileName := "test.log"
  logger, err := NewLogger(logDirPath, logFileName)
  if err != nil {
    t.Errorf("NewLogger(%s, %s) = %v", logDirPath, logFileName, err)
  }
  if logger == nil {
    t.Errorf("NewLogger(%s, %s) = %v", logDirPath, logFileName, logger)
  }
}
