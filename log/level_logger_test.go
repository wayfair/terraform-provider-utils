package log

import (
	"bytes"
	"log"
	"testing"
)

// -----------------------------------------------------------------------------
// Logging Level Constants / Enum Representation
// -----------------------------------------------------------------------------

// Ensures LogLevel.String() returns the proper string representation of
// the LogLevel (for valid LogLevel values).
func TestLogLevelStringValid(t *testing.T) {

	expectedLogLevelStrings := map[LogLevel]string{
		LevelDebug:   "DEBUG",
		LevelTrace:   "TRACE",
		LevelInfo:    "INFO",
		LevelWarning: "WARNING",
		LevelError:   "ERROR",
		LevelNone:    "NONE",
	}

	for key, value := range expectedLogLevelStrings {
		if key.String() != value {
			t.Errorf(
				"LogLevel.String did not return the correct string representation of "+
					"the LogLevel. Expected [%s], got [%s] for [%d].",
				value,
				key.String(),
				key,
			)
		}
	}

}

// Ensures LogLevel.String() returns the proper string representation of
// the LogLevel (for invalid LogLevel values).
func TestLogLevelStringInvalid(t *testing.T) {

	expectedLogLevelStrings := map[LogLevel]string{
		LogLevel(1000): "",
		LogLevel(-99):  "",
	}

	for key, value := range expectedLogLevelStrings {
		if key.String() != value {
			t.Errorf(
				"LogLevel.String did not return the correct string representation of "+
					"the LogLevel. Expected [%s], got [%s] for [%d].",
				value,
				key.String(),
				key,
			)
		}
	}

}

// Ensures LogLevelFromString() works with valid string representations (values
// returned by LogLevel.String())
func TestLogLevelFromStringValid(t *testing.T) {

	expectedLogLevel := map[string]LogLevel{
		"DEBUG":   LevelDebug,
		"TRACE":   LevelTrace,
		"INFO":    LevelInfo,
		"WARNING": LevelWarning,
		"ERROR":   LevelError,
		"NONE":    LevelNone,
	}

	for key, value := range expectedLogLevel {
		level, err := LogLevelFromString(key)
		if err != nil {
			t.Errorf(
				"LogLevelFromString returned an error for a valid string representation. "+
					"Expected [nil], got [%s] for [%s].",
				err.Error(),
				key,
			)
		} else if level != value {
			t.Errorf(
				"LogLevelFromString did not return the correct LogLevel for the given "+
					"string representation. Expected [%d], got [%d] for [%s].",
				value,
				level,
				key,
			)
		}
	}

}

// Ensures LogLevelFromString() is not case sensitive when dealing
// with the string representations.
func TestLogLevelFromStringCaseInsensitive(t *testing.T) {

	expectedLogLevel := map[string]LogLevel{
		"debug":   LevelDebug,
		"TrACE":   LevelTrace,
		"infO":    LevelInfo,
		"WaRnInG": LevelWarning,
		"eRrOr":   LevelError,
		"NONE":    LevelNone,
	}

	for key, value := range expectedLogLevel {
		level, err := LogLevelFromString(key)
		if err != nil {
			t.Errorf(
				"LogLevelFromString returned an error for a valid string representation. "+
					"Expected [nil], got [%s] for [%s].",
				err.Error(),
				key,
			)
		} else if level != value {
			t.Errorf(
				"LogLevelFromString did not return the correct LogLevel for the given "+
					"string representation. Expected [%d], got [%d] for [%s].",
				value,
				level,
				key,
			)
		}
	}

}

// Ensures that LogLevelFromString() returns an error when given bad input
// as well as LevelInvalid.
func TestLogLevelFromStringInvalid(t *testing.T) {

	expectedLogLevel := map[string]LogLevel{
		"  debug": LevelInvalid,
		"FOO":     LevelInvalid,
		"":        LevelInvalid,
		"  ":      LevelInvalid,
		"debugg":  LevelInvalid,
		"DEBUG ":  LevelInvalid,
	}

	for key, value := range expectedLogLevel {
		level, err := LogLevelFromString(key)
		if err == nil {
			t.Errorf(
				"LogLevelFromString did not return an error for bad string log level."+
					"Expected [error], got [nil] for [%s].",
				key,
			)
		} else if level != LevelInvalid {
			t.Errorf(
				"LogLevelFromString did not return the correct LogLevel for the given "+
					"string representation. Expected [%d], got [%d] for [%s].",
				value,
				level,
				key,
			)
		}
	}

}

// -----------------------------------------------------------------------------
// Provider Logger
// -----------------------------------------------------------------------------

// Helper function to assert the logger's buffer is empty
func assertBufferEmpty(t *testing.T, buffer *bytes.Buffer, loggerLevel LogLevel, messageLevel LogLevel) {
	messageStr := buffer.String()
	if messageStr != "" {
		t.Errorf(
			"Logger buffer is not empty when invoking a [%s] message print.  The "+
				"logger's level is set to [%s].  The logger should not have printed a "+
				"message to the buffer. Buffer contents: [%s].",
			messageLevel.String(),
			loggerLevel.String(),
			messageStr,
		)
	}
}

// Helper function to assert the logger's buffer is not empty.
func assertBufferNotEmpty(t *testing.T, buffer *bytes.Buffer, loggerLevel LogLevel, messageLevel LogLevel) {
	messageStr := buffer.String()
	if messageStr == "" {
		t.Errorf(
			"Logger buffer is empty when invoking a [%s] message print.  The "+
				"logger's level is set to [%s].  The logger should have printed a "+
				"message to the buffer. Buffer contents: [%s].",
			messageLevel.String(),
			loggerLevel.String(),
			messageStr,
		)
	}
}

// Ensures when the logger's level is set to 'LevelDebug', all logging functions
// of level debug or higher are printed.
func TestLogLevelDebugOutput(t *testing.T) {

	logLevel := LevelDebug

	// initialize the logger to write to a buffer
	buffer := new(bytes.Buffer)
	logger := NewLevelLogger(buffer, log.LstdFlags, logLevel)
	testMessage := "Hello, World!"

	// All log messages should be printed

	logger.Debugf(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelDebug)
	buffer.Reset()

	logger.Tracef(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelTrace)
	buffer.Reset()

	logger.Infof(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelInfo)
	buffer.Reset()

	logger.Warningf(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelWarning)
	buffer.Reset()

	logger.Errorf(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelError)
	buffer.Reset()
}

// Ensures when the logger's level is set to 'LevelTrace', all logging functions
// of level trace or higher are printed.
func TestLogLevelTraceOutput(t *testing.T) {

	logLevel := LevelTrace

	// initialize the logger to write to a buffer
	buffer := new(bytes.Buffer)
	logger := NewLevelLogger(buffer, log.LstdFlags, logLevel)
	testMessage := "Hello, World!"

	// Below Trace should not be printed

	logger.Debugf(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelDebug)
	buffer.Reset()

	// Trace and above should be printed

	logger.Tracef(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelTrace)
	buffer.Reset()

	logger.Infof(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelInfo)
	buffer.Reset()

	logger.Warningf(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelWarning)
	buffer.Reset()

	logger.Errorf(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelError)
	buffer.Reset()
}

// Ensures when the logger's level is set to 'LevelInfo', all logging functions
// of level info or higher are printed.
func TestLogLevelInfoOutput(t *testing.T) {

	logLevel := LevelInfo

	// initialize the logger to write to a buffer
	buffer := new(bytes.Buffer)
	logger := NewLevelLogger(buffer, log.LstdFlags, logLevel)
	testMessage := "Hello, World!"

	// Below Info should not be printed

	logger.Debugf(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelDebug)
	buffer.Reset()

	logger.Tracef(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelTrace)
	buffer.Reset()

	// Info and above should be printed

	logger.Infof(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelInfo)
	buffer.Reset()

	logger.Warningf(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelWarning)
	buffer.Reset()

	logger.Errorf(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelError)
	buffer.Reset()
}

// Ensures when the logger's level is set to 'LevelWarning', all logging functions
// of level warning or higher are printed.
func TestLogLevelWarningOutput(t *testing.T) {

	logLevel := LevelWarning

	// initialize the logger to write to a buffer
	buffer := new(bytes.Buffer)
	logger := NewLevelLogger(buffer, log.LstdFlags, logLevel)
	testMessage := "Hello, World!"

	// Below Warning should not be printed

	logger.Debugf(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelDebug)
	buffer.Reset()

	logger.Tracef(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelTrace)
	buffer.Reset()

	logger.Infof(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelInfo)
	buffer.Reset()

	// Warning and above should be printed

	logger.Warningf(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelWarning)
	buffer.Reset()

	logger.Errorf(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelError)
	buffer.Reset()
}

// Ensures when the logger's level is set to 'LevelError', all logging functions
// of level error or higher are printed.
func TestLogLevelErrorOutput(t *testing.T) {

	logLevel := LevelError

	// initialize the logger to write to a buffer
	buffer := new(bytes.Buffer)
	logger := NewLevelLogger(buffer, log.LstdFlags, logLevel)
	testMessage := "Hello, World!"

	// Anything under Error should not be printed

	logger.Debugf(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelDebug)
	buffer.Reset()

	logger.Tracef(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelTrace)
	buffer.Reset()

	logger.Infof(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelInfo)
	buffer.Reset()

	logger.Warningf(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelWarning)
	buffer.Reset()

	// Error and above should be printed

	logger.Errorf(testMessage)
	assertBufferNotEmpty(t, buffer, logLevel, LevelError)
	buffer.Reset()

}

// Ensures when the logger's level is set to 'LevelNone', no logging calls
// are written.
func TestLogLevelNoneOutput(t *testing.T) {

	logLevel := LevelNone

	// initialize the logger to write to a buffer
	buffer := new(bytes.Buffer)
	logger := NewLevelLogger(buffer, log.LstdFlags, logLevel)
	testMessage := "Hello, World!"

	// Nothing should be printed

	logger.Debugf(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelDebug)
	buffer.Reset()

	logger.Tracef(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelTrace)
	buffer.Reset()

	logger.Infof(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelInfo)
	buffer.Reset()

	logger.Warningf(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelWarning)
	buffer.Reset()

	logger.Errorf(testMessage)
	assertBufferEmpty(t, buffer, logLevel, LevelError)
	buffer.Reset()

}
