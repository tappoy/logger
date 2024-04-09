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

	// create logger
	logger := &Logger{
		logDirPath: logDirPath,
	}

	return logger, nil
}

func (logger *Logger) log(level string, now time.Time, format string, args ...interface{}) {
	// create log file if not exists
	if err := createFileIfNotExist(logger.logDirPath, level+".log"); err != nil {
		return
	}

	// get log file info
	logFilePath := filepath.Join(logger.logDirPath, level+".log")
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
		if _, err := os.Stat(filepath.Join(logger.logDirPath, rotateFileName)); err != nil {
			// if not exists, rename log file to rotate file
			os.Rename(logFilePath, filepath.Join(logger.logDirPath, rotateFileName))

			// Processing continues even if the rotation fails.
			// It is more fatal to fail to keep a log.
			// So, we don't check the error.

			if err := createFileIfNotExist(logger.logDirPath, level+".log"); err != nil {
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
	logFilePath := filepath.Join(logger.logDirPath, "debug.log")
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
