# About
This golang logging library provides these features:
- Logging to file.
- Log level filtering that contains: NONE, ERROR, WARN, INFO, DEBUG.

# Functions
- `NewLogger(logDirPath string, logFileName string) (*Logger, error)`: Create a new logger. Log level is set to INFO by default.
- `(*Logger) Error(format string, v ...interface{})`: Log error message.
- `(*Logger) Warn(format string, v ...interface{})`: Log warning message.
- `(*Logger) Info(format string, v ...interface{})`: Log info message.
- `(*Logger) Debug(format string, v ...interface{})`: Log debug message.
- `(*Logger) SetLogLevel(logLevel string) error`: Set log level.

# Errors
- `ErrInvalidLogLevel`: Invalid log level.
- `ErrCannotCreateLogDir`: Cannot create log directory.
- `ErrCannotCreateLogFile`: Cannot create log file.
- `ErrCannotWriteLogFile`: Cannot write log file.

# LICENSE
This module is licensed under the LGPL-3.0 license. For more information, see the LICENSE file.

# AUTHOR
[tappoy](https://github.com/tappoy)
