# About
This golang logging library provides these features:
- Logging to each level files. The files are below:
  - error.log: Error messages. Must be watched by the administrator.
  - notice.log: Messages that are not error but should be noted. Should be watched by the administrator.
  - info.log: Normal activity messages. Not necessary to be watched but helpful for the operation.
  - debug.log: Debug messages. For developers to debug. Should turn off in production.
- Debug output can be turned on if debug.log exists. if not exists, debug output is turned off.
- Log rotation. The log files are rotated when date changes. `ex) error.log -> error_2024-04-09.log`

# Format
Output logs in LTSV format.

LEVEL.log
```
datetime:YYYY-MM-DD HH:MM:SS\tLEVEL:LOG_MESSAGE\n
```

## Example of log
error.log
```
datetime:2024-04-05 20:37:04	error:message
```

notice.log
```
datetime:2024-04-05 20:37:04	notice:message
```

info.log
```
datetime:2024-04-05 20:37:04	info:message
```

debug.log
```
datetime:2024-04-05 20:37:04	debug:message
```

# Functions
- `NewLogger(logDirPath string) (*Logger, error)`: Create a new logger.
- `(*Logger) Error(format string, v ...interface{})`: Log error message to error.log.
- `(*Logger) Notice(format string, v ...interface{})`: Log notice message to notice.log.
- `(*Logger) Info(format string, v ...interface{})`: Log info message to info.log.
- `(*Logger) Debug(format string, v ...interface{})`: Log debug message to debug.log.

# Errors
- `ErrCannotCreateLogDir`: Cannot create log directory.
- `ErrCannotCreateLogFile`: Cannot create log file.
- `ErrCannotWriteLogFile`: Cannot write log file.

# LICENSE
[LGPL-3.0](LICENSE)

# AUTHOR
[tappoy](https://github.com/tappoy)
