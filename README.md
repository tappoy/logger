# About
This golang logging library provides these features:
- Logging to file.
- Filtering at the log level. Log levels are as follows:
  - DEBUG: Log all messages.
  - INFO: Log info and error messages.
  - ERROR: Log only error messages.
  - NONE: Do not log any messages.

# Format
Output logs in LTSV format.
```
datetime:YYYY-MM-DD HH:MM:SS\tlevel:LOG_LEVEL\tlog:LOG_MESSAGE\n
```

## Example of log
```
datetime:2024-04-05 20:37:04	level:DEBUG	log:debug message
datetime:2024-04-05 20:37:04	level:INFO	log:info message
datetime:2024-04-05 20:37:04	level:ERROR	log:error message
```

# Functions
- `NewLogger(logDirPath string, logFileName string, level string) (*Logger, error)`: Create a new logger.
- `(*Logger) Error(format string, v ...interface{})`: Log error message.
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
