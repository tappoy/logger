package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	logDirPath  string
	logFileName string
	level       string
}

var logLevelMap = map[string]int{
	"DEBUG": 0,
	"INFO":  1,
	"ERROR": 2,
	"NONE":  3,
}

func NewLogger(logDirPath string, logFileName string, level string) (*Logger, error) {
	// check log level
	if _, ok := logLevelMap[level]; !ok {
		return nil, errors.New("ErrInvalidLogLevel")
	}

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
		logDirPath:  logDirPath,
		logFileName: logFileName,
		level:       level,
	}

	return logger, nil
}

func (logger *Logger) SetLogLevel(level string) error {
	if _, ok := logLevelMap[level]; ok {
		logger.level = level
		return nil
	}

	return errors.New("ErrInvalidLogLevel")
}

func (logger *Logger) log(level string, format string, args ...interface{}) {
	// string to int log level
	levelInt := logLevelMap[level]
	limitInt := logLevelMap[logger.level]

	// filter log level
	if levelInt < limitInt {
		return
	}

	// log message
	logFilePath := filepath.Join(logger.logDirPath, logger.logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer logFile.Close()

	// log header, timestamp, level
	header := fmt.Sprintf("datetime:%s\tlevel:%s\tlog:", time.Now().Format("2006-01-02 15:04:05"), level)

	// Sprintf log message
	message := fmt.Sprintf(format, args...)

	// write log
	logFile.WriteString(header + message + "\n")

}

func (logger *Logger) Debug(format string, args ...interface{}) {
	logger.log("DEBUG", format, args...)
}

func (logger *Logger) Info(format string, args ...interface{}) {
	logger.log("INFO", format, args...)
}

func (logger *Logger) Error(format string, args ...interface{}) {
	logger.log("ERROR", format, args...)
}
