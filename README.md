# About
This golang logging library provides these features:
- Logging to file.
- Log level filtering that contains: NONE, ERROR, WARNING, INFO, DEBUG.

# Functions
- `NewLogger(logDir string, logFileName string, logLevel string) (*Logger, error)`: Create a new logger.
- `(*Logger) Error(format string, v ...interface{})`: Log error message.
- `(*Logger) Warning(format string, v ...interface{})`: Log warning message.
- `(*Logger) Info(format string, v ...interface{})`: Log info message.
- `(*Logger) Debug(format string, v ...interface{})`: Log debug message.
- `(*Logger) SetLogLevel(logLevel string) error`: Set log level.

# Errors
- `ErrInvalidLogLevel`: Invalid log level.
- `ErrCannotAccessLogDir`: Cannot access log directory.
- `ErrCannotAccessLogFile`: Cannot access log file.

# LICENSE
This module is licensed under the LGPL-3.0 license. For more information, see the LICENSE file.

# AUTHOR
[tappoy](https://github.com/tappoy)
