package logger

import (
  "os"
  "errors"
  "path/filepath"
)

type Logger struct {
  logDirPath string
  logFileName string
  logLevel string
}

var logLevels = []string{"DEBUG", "INFO", "WARN", "ERROR", "NONE"}

func NewLogger(logDirPath string, logFileName string) (*Logger, error) {
  // create log directory if not exists
  if _, err := os.Stat(logDirPath); os.IsNotExist(err) {
    err := os.MkdirAll(logDirPath, 0755)
    if err != nil {
      return nil, errors.New("ErrCannotCreateLogDir")
    }
  }

  // create log file if not exists
  logFilePath := filepath.Join(logDirPath, logFileName)
  if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
    _, err := os.Create(filepath.Join(logDirPath, logFileName))
    if err != nil {
      return nil, errors.New("ErrCannotCreateLogFile")
    }
  }

  // check log file write permission
  if _, err := os.OpenFile(logFilePath, os.O_WRONLY, 0666); err != nil {
    return nil, errors.New("ErrCannotWriteLogFile")
  }

  // create logger
  logger := &Logger{
    logDirPath: logDirPath,
    logFileName: logFileName,
    logLevel: "INFO",
  }

  return logger, nil
}
