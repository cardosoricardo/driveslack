package app

import (
	"log"
	"os"
	"runtime"
	"strings"
)

var logPath = "errors.log"

func checkError(err error) bool {
	if err != nil {
		logError(err)
		return true
	}
	return false
}

//logError returns the line on which the error was reported
func logError(e error) {
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(f)
		defer f.Close()
	}
	_, file, line, ok := runtime.Caller(2) // logger + checkError function.
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	log.Println(file, line, ":", e.Error())
}
