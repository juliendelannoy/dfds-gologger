package gologger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var ErrorCantCreateLogDirectory = errors.New("Can't create log directory")

func Init(logPath, complementaryInformation string) (err error) {
	initLogger(os.Stderr)
	mw, err := setLogFile(logPath, complementaryInformation)
	if err != nil {
		return err
	}
	initLogger(mw)
	return nil
}

func initLogger(
	output io.Writer) {

	Debug = log.New(output,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lmicroseconds) //|log.Lshortfile)

	Info = log.New(output,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lmicroseconds)

	Warning = log.New(output,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lmicroseconds)

	Error = log.New(output,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

}

func getLogFilename(complementaryInformation string) (s string) {
	appName := filepath.Base(os.Args[0])
	logFilename := fmt.Sprintf("%s_%s_%s_pid%d.log", time.Now().Format("20060102T150405"),
		appName, complementaryInformation, os.Getpid())

	return logFilename
}

func createLogDirectory(logPath string) (err error) {
	err = os.Mkdir(logPath, 0744)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return ErrorCantCreateLogDirectory
	}
}

func newMultiWriter(logFileName string) (w io.Writer, err error) {

	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return w, err
	}
	return io.MultiWriter(os.Stdout, file), nil
}

func setLogFile(logPath, complementaryInformation string) (w io.Writer, err error) {
	err = createLogDirectory(logPath)
	if err != nil {
		return nil, err
	}
	logFileName := getLogFilename(complementaryInformation)
	fullLogFileName := filepath.Join(logPath, logFileName)
	return newMultiWriter(fullLogFileName)
}
