package logger

import (
	"log"
	"os"
)

var infoLog Logger
var errorLog Logger

type Logger interface {
	log() *log.Logger
}
type Infolog struct {
	Infolog *log.Logger
}

type Errorlog struct {
	Errorlog *log.Logger
}

func (i Infolog) log() *log.Logger {

	log := i.Infolog
	return log
}

func (e Errorlog) log() *log.Logger {

	log := e.Errorlog
	return log
}

func Info(format string, a ...interface{}) {

	infoLog.log().Printf(format, a...)
}

func Infoln(a ...interface{}) {

	infoLog.log().Println(a...)
}

func Infofatal(a ...interface{}) {

	infoLog.log().Fatal(a...)
}

func Error(format string, a ...interface{}) {

	errorLog.log().Printf(format, a...)
}

func Errorln(a ...interface{}) {

	errorLog.log().Println(a...)
}

func Errorfatal(a ...interface{}) {

	errorLog.log().Fatal(a...)
}

func init() {

	infoLog = &Infolog{log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)}
	errorLog = &Errorlog{log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime)}

}
