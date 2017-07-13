package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
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

//Save insert relation folderID with last date get informacion in the file
func Save(folderID string, lastDate time.Time) {
	data := map[string]time.Time{}
	pathDB := ("db.txt")
	// s.OpenFile(pathDB, os.O_RDONLY|os.O_CREATE, 0666)
	if _, err := os.Stat(pathDB); os.IsNotExist(err) {
		_, err := os.Create(pathDB)
		if checkError(err) {
			return
		}
	}

	file, err := os.OpenFile(pathDB, os.O_RDWR, 0644)
	if checkError(err) {
		return
	}
	defer file.Close()

	if f, _ := file.Stat(); f.Size() > 0 {
		err = json.NewDecoder(file).Decode(&data)
		if checkError(err) {
			return
		}
	}

	if data[folderID] == lastDate {
		return
	}

	err = file.Truncate(0) //empty file
	if checkError(err) {
		return
	}

	_, err = file.Seek(0, 0)
	if checkError(err) {
		return
	}

	data[folderID] = lastDate

	err = json.NewEncoder(file).Encode(data)
	if checkError(err) {
		return
	}

	err = file.Sync()
	if checkError(err) {
		return
	}
}

//Get obtains relation folderID with last date get informacion
func Get(folderID string) (lastDate time.Time) {
	data := map[string]time.Time{}
	pathDB := ("db.txt")

	if _, err := os.Stat(pathDB); os.IsNotExist(err) {
		_, err := os.Create(pathDB)
		if checkError(err) {
			return
		}
	}

	file, err := os.OpenFile(pathDB, os.O_RDWR, 0644)
	if checkError(err) {
		return
	}

	if f, ex := file.Stat(); f.Size() > 0 {
		if checkError(ex) {
			return
		}
		err = json.NewDecoder(file).Decode(&data)
		if checkError(err) {
			return
		}
	}

	defer file.Close()
	if data[folderID] != (time.Time{}) {
		lastDate = data[folderID]
	}
	return
}
