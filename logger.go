// This golang package provides these features:
//
// Logging to each level files. The files are below:
//
//   - error.log: Error messages. Must be watched by the administrator.
//   - notice.log: Messages that are not error but should be noted. Should be watched by the administrator.
//   - info.log: Normal activity messages. Not necessary to be watched but helpful for the operation.
//   - debug.log: Debug messages. For developers to debug. Should turn off in production.
//   - logger.log: Logger's own fatal messages. Should be watched by the administrator.
//
// Debug output can be turned on if debug.log exists. if not exists, debug output is turned off.
//
// Log rotation. The log files are rotated when date changes.
//
//	ex) error.log -> backup/error_2024-04-09.log
//
// If there is over 30 files in backup directory, the oldest file is deleted.
//
// Output logs in LTSV format.
//
//	datetime:YYYY-MM-DD HH:MM:SS\tLEVEL:LOG_MESSAGE\n
//
// Output example:
//
//	datetime:2024-04-05 20:37:04	error:message    // error.log
//	datetime:2024-04-05 20:37:04	notice:message   // notice.log
//	datetime:2024-04-05 20:37:04	info:message     // info.log
//	datetime:2024-04-05 20:37:04	debug:message    // debug.log
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

var (
	// Cannot create log directory.
	ErrCannotCreateLogDir = errors.New("ErrCannotCreateLogDir")

	// Cannot create log file.
	ErrCannotCreateLogFile = errors.New("ErrCannotCreateLogFile")

	// Cannot write log file.
	ErrCannotWriteLogFile = errors.New("ErrCannotWriteLogFile")
)

func createFileIfNotExist(logFilePath string) error {
	if _, err := os.Stat(logFilePath); err != nil {
		_, err := os.Create(logFilePath)
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

// Create new logger.
//
// Errors:
//   - ErrCannotCreateLogDir
//   - ErrCannotWriteLogFile
func NewLogger(logDir string) (*Logger, error) {
	// create log directory if not exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0775)
		if err != nil {
			return nil, ErrCannotCreateLogDir
		}
	}

	// create tempfile to check write permission
	if tempfile, err := os.CreateTemp(logDir, "tempfile"); err != nil {
		return nil, ErrCannotWriteLogFile
	} else {
		os.Remove(tempfile.Name())
	}

	// create logger
	logger := &Logger{
		logDir: logDir,
	}

	return logger, nil
}

// rotate log file if necessary.
func (logger *Logger) rotate(logFilePath string, now time.Time, level string) error {
	// get log file info
	stat, err := os.Stat(logFilePath)
	if err != nil {
		// no file, no rotate
		return nil
	}

	// get modified time
	modTime := stat.ModTime()

	if modTime.Day() != now.Day() {
		// create backup directory if not exists
		if _, err := os.Stat(filepath.Join(logger.logDir, "backup")); err != nil {
			err := os.MkdirAll(filepath.Join(logger.logDir, "backup"), 0775)
			if err != nil {
				return err
			}
		}

		// remove old backup files
		list, err := filepath.Glob(filepath.Join(logger.logDir, "backup", "*"))
		if err != nil {
			return err
		}

		// check if backup files count is over 30
		if len(list) > 30 {
			// sort by modified time
			for i := 0; i < len(list); i++ {
				for j := i + 1; j < len(list); j++ {
					stat1, _ := os.Stat(list[i])
					stat2, _ := os.Stat(list[j])
					if stat1.ModTime().Before(stat2.ModTime()) {
						list[i], list[j] = list[j], list[i]
					}
				}
			}
			// remove old backup files over 30
			for i := 29; i < len(list); i++ {
				os.Remove(list[i])
			}
		}

		// make rote file name
		rotateFileName := level + "_" + modTime.Format("2006-01-02") + ".log"
		rotateFilePath := filepath.Join(logger.logDir, "backup", rotateFileName)
		// check if rotate file exists
		if _, err := os.Stat(rotateFilePath); err != nil {
			// if not exists, rename log file to rotate file
			os.Rename(logFilePath, rotateFilePath)
		}
	}

	return nil
}

func write(logFile *os.File, level string, now time.Time, format string, args ...any) {
	// log header, timestamp, level
	header := fmt.Sprintf("datetime:%s\t%s:", now.Format("2006-01-02 15:04:05"), level)

	// Sprintf log message
	message := fmt.Sprintf(format, args...)

	// write log
	logFile.WriteString(header + message + "\n")
}

func (logger *Logger) fatal(now time.Time, message string) {
	logFilePath := filepath.Join(logger.logDir, "logger.log")
	err := createFileIfNotExist(logFilePath)
	if err != nil {
		return
	}

	// open log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer logFile.Close()

	write(logFile, "logger", now, "%s", message)
}

func (logger *Logger) log(level string, now time.Time, format string, args ...any) {
	logFilePath := filepath.Join(logger.logDir, level+".log")

	// rotate if log file is not today's.
	err := logger.rotate(logFilePath, now, level)
	if err != nil {
		logger.fatal(now, "rotate error: "+err.Error())
		// Processing continues even if the rotation fails.
		// It is more fatal to fail to keep a log.
	}

	// create log file if not exists
	err = createFileIfNotExist(logFilePath)
	if err != nil {
		logger.fatal(now, "create log file error: "+err.Error())
		return
	}

	// open log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logger.fatal(now, "open file error: "+err.Error())
		return
	}
	defer logFile.Close()

	write(logFile, level, now, format, args...)
}

// Log debug message.
func (logger *Logger) Debug(format string, args ...any) {
	// check if debug log file exists
	logFilePath := filepath.Join(logger.logDir, "debug.log")
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		return
	}
	logger.log("debug", time.Now(), format, args...)
}

// Log info message.
func (logger *Logger) Info(format string, args ...any) {
	logger.log("info", time.Now(), format, args...)
}

// Log notice message.
func (logger *Logger) Notice(format string, args ...any) {
	logger.log("notice", time.Now(), format, args...)
}

// Log error message.
func (logger *Logger) Error(format string, args ...any) {
	logger.log("error", time.Now(), format, args...)
}

// Get log directory.
func (logger *Logger) GetLogDir() string {
	return logger.logDir
}
