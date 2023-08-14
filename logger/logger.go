package logger

import (
	"io"
	"log"
	"os"

	"github.com/minhnvc/go_corelib/utils"
)

var isLog bool
var isLogDB bool
var isDev bool

func InitLogger() {
	//init log
	pathLog := "/data/log/zaloanalytics-v2/zaloanalytics-v2.log"
	if utils.GetConfig("ROOT_DOMAIN") == "localhost" {
		isDev = true
		pathLog = "za_analytics.log"
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
	if !isDev && category == "dev" {
		return
	}
	if isLog {
		if category == "Mongo" && !isLogDB {
			return
		}
		log.Println("["+category+"]", d)
	}
}
