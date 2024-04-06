# About
This golang logging library provides these features:
- Logging to each level files (error.log, info.log, debug.log).
- Debug output can be turned on/off.

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

info.log
```
datetime:2024-04-05 20:37:04	info:message
```

debug.log
```
datetime:2024-04-05 20:37:04	debug:message
```

# Functions
- `NewLogger(logDirPath string, debug bool) (*Logger, error)`: Create a new logger.
- `(*Logger) Error(format string, v ...interface{})`: Log error message to error.log.
- `(*Logger) Info(format string, v ...interface{})`: Log info message to info.log.
- `(*Logger) Debug(format string, v ...interface{})`: Log debug message to debug.log.
- `(*Logger) SetDebug(debug bool)`: Set debug mode on/off.

# Errors
- `ErrCannotCreateLogDir`: Cannot create log directory.
- `ErrCannotCreateLogFile`: Cannot create log file.
- `ErrCannotWriteLogFile`: Cannot write log file.

# LICENSE
This module is licensed under the LGPL-3.0 license. For more information, see the LICENSE file.

# AUTHOR
[tappoy](https://github.com/tappoy)
