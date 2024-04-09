package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	logDirPath string
}

func createFileIfNotExist(dirPath string, fileName string) error {
	// create log file if not exists
	logFilePath := filepath.Join(dirPath, fileName)
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		_, err := os.Create(filepath.Join(dirPath, fileName))
		if err != nil {
			return errors.New("ErrCannotCreateLogFile")
		}
	}

	// check log file write permission
	if _, err := os.OpenFile(logFilePath, os.O_WRONLY, 0666); err != nil {
		return errors.New("ErrCannotWriteLogFile")
	}

	return nil
}

func NewLogger(logDirPath string) (*Logger, error) {
	// create log directory if not exists
	if _, err := os.Stat(logDirPath); os.IsNotExist(err) {
		err := os.MkdirAll(logDirPath, 0755)
		if err != nil {
			return nil, errors.New("ErrCannotCreateLogDir")
		}
	}

	// create error log file if not exists
	err := createFileIfNotExist(logDirPath, "error.log")
	if err != nil {
		return nil, err
	}

	// create notice log file if not exists
	err = createFileIfNotExist(logDirPath, "notice.log")
	if err != nil {
		return nil, err
	}

	// create info log file if not exists
	err = createFileIfNotExist(logDirPath, "info.log")
	if err != nil {
		return nil, err
	}

	// create logger
	logger := &Logger{
		logDirPath: logDirPath,
	}

	return logger, nil
}

func (logger *Logger) log(level string, format string, args ...interface{}) {
	// open log file
	logFilePath := filepath.Join(logger.logDirPath, level+".log")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer logFile.Close()

	// log header, timestamp, level
	header := fmt.Sprintf("datetime:%s\t%s:", time.Now().Format("2006-01-02 15:04:05"), level)

	// Sprintf log message
	message := fmt.Sprintf(format, args...)

	// write log
	logFile.WriteString(header + message + "\n")
}

func (logger *Logger) Debug(format string, args ...interface{}) {
	// check if debug log file exists
	logFilePath := filepath.Join(logger.logDirPath, "debug.log")
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		return
	}
	logger.log("debug", format, args...)
}

func (logger *Logger) Info(format string, args ...interface{}) {
	logger.log("info", format, args...)
}

func (logger *Logger) Notice(format string, args ...interface{}) {
	logger.log("notice", format, args...)
}

func (logger *Logger) Error(format string, args ...interface{}) {
	logger.log("error", format, args...)
}
