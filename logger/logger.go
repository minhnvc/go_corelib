package logger

import (
	"io"
	"log"
	"os"

	"github.com/minhnvc/go_corelib/utils"
)

var isLog bool
var isLogDB bool

func InitLogger(name string) {
	//init log
	pathLog := "/data/log/" + name + "/" + name + ".log"
	if utils.GetConfig("IS_LOCAL") == "true" {
		pathLog = name + ".log"
	}
	logFile, err := os.OpenFile(pathLog, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	isLog = utils.GetConfig("LOG_CONSOLE") == "true"
	isLogDB = utils.GetConfig("LOG_DB_CONSOLE") == "true"
}

func PrintLn(category string, d ...interface{}) {
	if isLog {
		if category == "Mongo" && !isLogDB {
			return
		}
		log.Println("["+category+"]", d)
	}
}
