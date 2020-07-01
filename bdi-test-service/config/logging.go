//Package config ... custom loggers
package config

import (
	"log"
	"os"
)

//Logging ...
type Logging struct {
	Infolog  *log.Logger
	Errorlog *log.Logger
}

//Logger ... custom logger
func Logger() *Logging {

	log := &Logging{
		Infolog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		Errorlog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)}

	return log
}
