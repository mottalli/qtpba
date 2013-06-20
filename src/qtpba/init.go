package qtpba

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	initLogger()
	initDB()
	initBlacklist()
}

func initLogger() {
	logger = log.New(os.Stdout, "qtpba - ", log.Ldate|log.Ltime)
}
