package main

import (
	"encoding/json"
	"io/ioutil"
)

// GetRelationFromFile obtains the relation of drive and slack from file
func GetRelationFromFile(path string) (ds []DriveSlack) {
	file, err := ioutil.ReadFile(path)
	if checkError(err) {
		return
	}
	json.Unmarshal(file, &ds)
	return
}
