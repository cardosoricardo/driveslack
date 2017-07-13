package main

import (
	"fmt"
	"runtime"
	"strings"
)

func checkError(err error) bool {
	if err != nil {
		LogError(err)
		return true
	}
	return false
}

//LogError returns the line on which the error was reported
func LogError(e error) {
	_, file, line, ok := runtime.Caller(1) // logger + public function.
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
	fmt.Println("log", file, line)
}
