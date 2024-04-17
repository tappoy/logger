package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	logDir string
}

// Errors
var (
  ErrCannotCreateLogDir = errors.New("ErrCannotCreateLogDir")
  ErrCannotCreateLogFile = errors.New("ErrCannotCreateLogFile")
  ErrCannotWriteLogFile = errors.New("ErrCannotWriteLogFile")
)

func createFileIfNotExist(dirPath string, fileName string) error {
	// create log file if not exists
	logFilePath := filepath.Join(dirPath, fileName)
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		_, err := os.Create(filepath.Join(dirPath, fileName))
		if err != nil {
			return ErrCannotCreateLogFile
		}
	}

	// check log file write permission
	if _, err := os.OpenFile(logFilePath, os.O_WRONLY, 0666); err != nil {
		return ErrCannotWriteLogFile
	}

	return nil
}

func NewLogger(logDir string) (*Logger, error) {
	// create log directory if not exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			return nil, ErrCannotCreateLogDir
		}
	}

	// create logger
	logger := &Logger{
		logDir: logDir,
	}

	return logger, nil
}

func (logger *Logger) log(level string, now time.Time, format string, args ...interface{}) {
	// create log file if not exists
	if err := createFileIfNotExist(logger.logDir, level+".log"); err != nil {
		return
	}

	// get log file info
	logFilePath := filepath.Join(logger.logDir, level+".log")
	stat, err := os.Stat(logFilePath)
	if err != nil {
		return
	}

	// get modified time
	modTime := stat.ModTime()

	// rotate if log file is not today's
	if modTime.Day() != now.Day() {
		// make rote file name
		rotateFileName := level + "_" + modTime.Format("2006-01-02") + ".log"
		// check if rotate file exists
		if _, err := os.Stat(filepath.Join(logger.logDir, rotateFileName)); err != nil {
			// if not exists, rename log file to rotate file
			os.Rename(logFilePath, filepath.Join(logger.logDir, rotateFileName))

			// Processing continues even if the rotation fails.
			// It is more fatal to fail to keep a log.
			// So, we don't check the error.

			if err := createFileIfNotExist(logger.logDir, level+".log"); err != nil {
				return
			}
		}
	}

	// open log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer logFile.Close()

	// log header, timestamp, level
	header := fmt.Sprintf("datetime:%s\t%s:", now.Format("2006-01-02 15:04:05"), level)

	// Sprintf log message
	message := fmt.Sprintf(format, args...)

	// write log
	logFile.WriteString(header + message + "\n")
}

func (logger *Logger) Debug(format string, args ...interface{}) {
	// check if debug log file exists
	logFilePath := filepath.Join(logger.logDir, "debug.log")
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		return
	}
	logger.log("debug", time.Now(), format, args...)
}

func (logger *Logger) Info(format string, args ...interface{}) {
	logger.log("info", time.Now(), format, args...)
}

func (logger *Logger) Notice(format string, args ...interface{}) {
	logger.log("notice", time.Now(), format, args...)
}

func (logger *Logger) Error(format string, args ...interface{}) {
	logger.log("error", time.Now(), format, args...)
}

func (logger *Logger) GetLogDir() string {
  return logger.logDir
}
