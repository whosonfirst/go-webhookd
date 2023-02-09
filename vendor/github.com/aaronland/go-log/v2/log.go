package log

import (
	"fmt"
	"log"
	"os"
)

// DEBUG_LEVEL is the numeric level associated with "debug" messages.
const DEBUG_LEVEL int = 10

// INFO_LEVEL is the numeric level associated with "info" messages.
const INFO_LEVEL int = 20

// WARNING_LEVEL is the numeric level associated with "warning" messages.
const WARNING_LEVEL int = 30

// ERROR_LEVEL is the numeric level associated with "error" messages.
const ERROR_LEVEL int = 40

// FATAL_LEVEL is the numeric level associated with "fatal" messages.
const FATAL_LEVEL int = 50

// DEBUG_PREFIX is the string prefix to prepend "debug" messages with.
const DEBUG_PREFIX string = "🪵" // "🦜"

// INFO_PREFIX is the string prefix to prepend "info" messages with.
const INFO_PREFIX string = "💬"

// WARNING_PREFIX is the string prefix to prepend "warning" messages with.
const WARNING_PREFIX string = "🧯"

// ERROR_PREFIX is the string prefix to prepend "error" messages with.
const ERROR_PREFIX string = "🔥"

// FATAL_PREFIX is the string prefix to prepend "fatal" messages with.
const FATAL_PREFIX string = "💥"

var minLevel int

// UnsetMinLevel resets the minimum level log level allowing all messages to be emitted.
func UnsetMinLevel() {
	minLevel = 0
}

// SetMinLevel assigns the minimum log level to 'l'. Log events with a lower level will not be emitted.
func SetMinLevel(l int) error {

	switch l {
	case DEBUG_LEVEL, INFO_LEVEL, WARNING_LEVEL, ERROR_LEVEL, FATAL_LEVEL:
		minLevel = l
	default:
		return fmt.Errorf("Invalid level")
	}

	return nil
}

// SetMinLevelStringWithPrefix assigns the minimum log level associated with 'prefix'. Log events with a lower level will not be emitted.
func SetMinLevelWithPrefix(prefix string) error {

	switch prefix {
	case DEBUG_PREFIX:
		return SetMinLevel(DEBUG_LEVEL)
	case INFO_PREFIX:
		return SetMinLevel(INFO_LEVEL)
	case WARNING_PREFIX:
		return SetMinLevel(WARNING_LEVEL)
	case ERROR_PREFIX:
		return SetMinLevel(ERROR_LEVEL)
	case FATAL_PREFIX:
		return SetMinLevel(FATAL_LEVEL)
	default:
		return fmt.Errorf("Invalid level")
	}
}

func emit(prefix string, args ...interface{}) {

	count_args := len(args)

	var logger *log.Logger
	var msg string
	var extras []interface{}

	// Nothing to do. Go home.

	if count_args == 0 {
		return
	}

	// Check to see whether first arg is a *log.Logger instance
	// If not create a logger and check whether first argument
	// is an error.

	if count_args >= 1 {

		switch args[0].(type) {
		case *log.Logger:
			logger = args[0].(*log.Logger)
		default:

			// See the way we are calling log.New rather than log.Default?
			// That's mostly so for the tests so that we can capture STDERR
			// The value if log.Default is a global singleton in log/log.go
			// which references an instance os.Stderr which doesn't seem to
			// get updated when we reassign it in the tests. I suppose this
			// might be the cause of "hilarity" in the future but it will
			// do for now...
			// https://cs.opensource.google/go/go/+/refs/tags/go1.20:src/log/log.go;l=89

			logger = log.New(os.Stderr, "", log.LstdFlags)

			switch args[0].(type) {
			case string:
				msg = args[0].(string)
			default:
				msg = fmt.Sprintf("%v", args[0])
			}
		}
	}

	// Check to see whether second arg is a log formatting string
	// or an error (or really anything other than a string)

	if count_args >= 2 {

		switch args[1].(type) {
		case string:
			msg = args[1].(string)
		default:
			msg = fmt.Sprintf("%v", args[1])
		}
	}

	// Anything else

	if count_args >= 3 {
		extras = args[2:]
	}

	msg = fmt.Sprintf("%s %s", prefix, msg)
	logger.Printf(msg, extras...)
}

// Emit a "debug" log message.
func Debug(args ...interface{}) {

	if minLevel > DEBUG_LEVEL {
		return
	}

	emit(DEBUG_PREFIX, args...)
}

// Emit a "info" log message.
func Info(args ...interface{}) {

	if minLevel > INFO_LEVEL {
		return
	}

	emit(INFO_PREFIX, args...)
}

// Emit a "warning" log message.
func Warning(args ...interface{}) {

	if minLevel > WARNING_LEVEL {
		return
	}

	emit(WARNING_PREFIX, args...)
}

// Emit an "error" log message.
func Error(args ...interface{}) {

	if minLevel > ERROR_LEVEL {
		return
	}

	emit(ERROR_PREFIX, args...)
}

// Emit an "fatal" log message.
func Fatal(args ...interface{}) {

	if minLevel > FATAL_LEVEL {
		return
	}

	emit(FATAL_PREFIX, args...)
	os.Exit(1)
}
