//logging package is a logging wrapper for the zero-log package.

package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger *zerolog.Logger
}

//New will instantiate a new instance of *Logger. Use this if printing logs
//to JSON outputs.
func New(output io.Writer, isDebug bool) *Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	//output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return &Logger{logger: &logger}

}

//NewConsole is similar to new but will print to os.Stdout by default with
//terminal formatting.
func NewConsole(isDebug bool) *Logger {

	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.WarnLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return &Logger{logger: &logger}
}

func logFileWriter() *os.File {
	now := time.Now()
	logFileSuffix := fmt.Sprintf("%d_%d_%d", now.Year(), now.Day(), now.Month())
	logFileName := "logfile" + logFileSuffix + ".log"

	f, err := os.OpenFile(logFileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	//TODO: defer close when shutting down main-program
	return f
}

//NewLogMultiWriterWRter will allow log output to multiple writers that implement the Writer interface
func NewLogMultiWriter(isDebug bool, writers ...io.Writer) *Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	multi := zerolog.MultiLevelWriter(writers...)

	logger := zerolog.New(multi).With().Timestamp().Logger()
	return &Logger{logger: &logger}
}

// Output duplicates the global logger and sets w as its output.
func (l *Logger) Output(w io.Writer) zerolog.Logger {
	return l.logger.Output(w)
}

// With creates a child logger with the field added to its context.
func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func (l *Logger) Level(level zerolog.Level) zerolog.Logger {
	return l.logger.Level(level)
}

// Sample returns a logger with the s sampler.
func (l *Logger) Sample(s zerolog.Sampler) zerolog.Logger {
	return l.logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func (l *Logger) Hook(h zerolog.Hook) zerolog.Logger {
	return l.logger.Hook(h)
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

//Infoprint is a shortcut to printing out errors regardless of whether
//debug level is set to true or false.
//The usual way with zerolog for command chaining for errors and info is as follows:
//log.Info().Msg(). Furthermore, the default .Msg() does not take an interface as an
//argument. So type error messages can be passed to below.
//So below allows for printing errors without the need to chain if debug is disabled
func (l *Logger) Infoprint(v ...interface{}) {
	e := l.Info()
	e.Msg(fmt.Sprint(v...))
}

//Infoprintf see comment for above
func (l *Logger) Infoprintf(format string, v ...interface{}) {
	e := l.Info()
	e.Msgf(format, v...)
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

//PrintWarn ... see comment for PrintInfo
func (l *Logger) PrintWarn(v ...interface{}) {
	e := l.Warn()
	e.Msg(fmt.Sprint(v...))
}

//PrintWarnF see comment for above
func (l *Logger) PrintWarnF(format string, v ...interface{}) {
	e := l.Warn()
	e.Msgf(format, v...)
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

//PrintError see comment for PrintInfo
func (l *Logger) PrintError(v ...interface{}) {
	e := l.Error()
	e.Msg(fmt.Sprint(v...))
}

//PrintErrorf see comment for PrintInfo
func (l *Logger) PrintErrorf(format string, v ...interface{}) {
	e := l.Error()
	e.Msgf(format, v...)
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

//PrintFatal calls Fatal with os.Exit(1) and prints the error
func (l *Logger) PrintFatal(v ...interface{}) {
	e := l.Fatal()
	e.Msg(fmt.Sprint(v...))
}

//PrintFatalf calls Fatal with os.Exit(1) and prints the error
func (l *Logger) PrintFatalf(format string, v ...interface{}) {
	e := l.Fatal()
	e.Msgf(format, v...)
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Log() *zerolog.Event {
	return l.logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	l.logger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func (l *Logger) Ctx(ctx context.Context) *Logger {
	return &Logger{logger: zerolog.Ctx(ctx)}
}
