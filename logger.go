package datalogger

import (
	"datalogger/src/file"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
)

// LogLevel type
type LogLevel int32

const (
	//TRACE logs everything
	TRACE LogLevel = 1 << iota
	//DEBUG logs Debug, Info, Warnings and Errors
	DEBUG
	//INFO logs Info, Warnings and Errors
	INFO
	//WARN logs Warning and Errors
	WARN
	//ERROR logs just Errors
	ERROR
)

// DefaultFlags used by created loggers
var DefaultFlags = log.Ldate | log.Ltime | log.Lshortfile

type configLogger struct {
	loglevel int32
	flags    int

	ptrLoggerDbg   *log.Logger
	ptrLoggerTrace *log.Logger
	ptrLoggerInfo  *log.Logger
	ptrLoggerWarn  *log.Logger
	ptrLoggerError *log.Logger

	LogFile *file.RotatingFileLog
}

var Logger configLogger

func doLogging(logLevel LogLevel, fileName string, maxBytes, backupCount int) {
	traceHandle := ioutil.Discard
	debugHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard

	var fileHandle *file.RotatingFileLog

	switch logLevel {
	case TRACE:
		traceHandle = os.Stdout
		fallthrough
	case DEBUG:
		debugHandle = os.Stdout
		fallthrough
	case INFO:
		infoHandle = os.Stdout
		fallthrough
	case WARN:
		warnHandle = os.Stdout
		fallthrough
	case ERROR:
		errorHandle = os.Stderr
	}

	if fileName != "" {
		var err error
		fileHandle, err = file.NewRotatingFileLog(fileName, maxBytes, backupCount)
		if err != nil {
			log.Fatal("logger: unable to create RotatingFileHandler: ", err)
		}

		if traceHandle == os.Stdout {
			traceHandle = io.MultiWriter(fileHandle, traceHandle)
		}
		if debugHandle == os.Stdout {
			debugHandle = io.MultiWriter(fileHandle, traceHandle)
		}

		if infoHandle == os.Stdout {
			infoHandle = io.MultiWriter(fileHandle, infoHandle)
		}

		if warnHandle == os.Stdout {
			warnHandle = io.MultiWriter(fileHandle, warnHandle)
		}

		if errorHandle == os.Stderr {
			errorHandle = io.MultiWriter(fileHandle, errorHandle)
		}
	}

	Logger.flags = DefaultFlags
	Logger = configLogger{
		ptrLoggerTrace: log.New(traceHandle, "T: ", Logger.flags),
		ptrLoggerDbg:   log.New(debugHandle, "D: ", Logger.flags),
		ptrLoggerInfo:  log.New(infoHandle, "I: ", Logger.flags),
		ptrLoggerWarn:  log.New(warnHandle, "W: ", Logger.flags),
		ptrLoggerError: log.New(errorHandle, "E: ", Logger.flags),
		LogFile:        fileHandle,
	}
	atomic.StoreInt32(&Logger.loglevel, int32(logLevel))
}

// Start starts the logging
func Start(level LogLevel, path string) {
	doLogging(level, path, file.MaxSizeFile, file.BckCount)
}

func StartEx(level LogLevel, path string, maxBytes, backupCount int) {
	doLogging(level, path, maxBytes, backupCount)
}

// Stop stops the logging
func Stop() error {
	if Logger.LogFile != nil {
		return Logger.LogFile.Close()
	}
	return nil
}

//Sync commits the current contents of the file to stable storage.
//Typically, this means flushing the file system's in-memory copy
//of recently written data to disk.
func Sync() {
	if Logger.LogFile != nil {
		Logger.LogFile.Fdescr.Sync()
	}
}

// Trace writes to the Trace destination
func Trace(format string, a ...interface{}) {
	Logger.ptrLoggerTrace.Output(2, fmt.Sprintf(format, a...))
}

// Debug writes to the Debug destination
func Debug(format string, a ...interface{}) {
	Logger.ptrLoggerDbg.Output(2, fmt.Sprintf(format, a...))
}

// Info writes to the Info destination
func Info(format string, a ...interface{}) {
	Logger.ptrLoggerInfo.Output(2, fmt.Sprintf(format, a...))
}

// Warning writes to the Warning destination
func Warning(format string, a ...interface{}) {
	Logger.ptrLoggerWarn.Output(2, fmt.Sprintf(format, a...))
}

// Error writes to the Error destination and accepts an err
func Error(err error) {
	Logger.ptrLoggerError.Output(2, fmt.Sprintf("%s\n", err))
}

// IfError is a shortcut function for log.Error if error
func IfError(err error) {
	if err != nil {
		Logger.ptrLoggerError.Output(2, fmt.Sprintf("%s\n", err))
	}
}
