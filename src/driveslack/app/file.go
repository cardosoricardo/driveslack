package app

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

var dbPath = "db.txt"

//DriveSlack defines a struct of relation between drive id and slack id
type DriveSlack struct {
	ChannelID string `json:"channel"`
	DriveID   string `json:"drive"`
}

// GetRelationFromFile obtains the relation of drive and slack from file
func GetRelationFromFile(path string) (ds []DriveSlack, err error) {
	file, err := ioutil.ReadFile(path)
	if checkError(err) {
		return
	}
	err = json.Unmarshal(file, &ds)
	checkError(err)
	return
}

//Save insert relation folderID with last date get informacion in the file
func Save(folderID string, lastDate time.Time) error {
	data := map[string]time.Time{}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		_, err := os.Create(dbPath)
		if checkError(err) {
			return err
		}
	}

	file, err := os.OpenFile(dbPath, os.O_RDWR, 0644)
	if checkError(err) {
		return err
	}
	defer file.Close()

	if f, _ := file.Stat(); f.Size() > 0 {
		err = json.NewDecoder(file).Decode(&data)
		if checkError(err) {
			return err
		}
	}

	if data[folderID] == lastDate {
		return nil
	}

	err = file.Truncate(0) //empty file
	if checkError(err) {
		return err
	}

	_, err = file.Seek(0, 0)
	if checkError(err) {
		return err
	}

	data[folderID] = lastDate

	err = json.NewEncoder(file).Encode(data)
	if checkError(err) {
		return err
	}

	err = file.Sync()
	checkError(err)
	return err
}

//get obtains relation folderID with last date get informacion
func get(folderID string) (lastDate time.Time) {
	data := map[string]time.Time{}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return
	}

	file, err := os.OpenFile(dbPath, os.O_RDWR, 0644)
	if checkError(err) {
		return
	}
	defer file.Close()

	if f, ex := file.Stat(); f.Size() > 0 {
		if checkError(ex) {
			return
		}
		err = json.NewDecoder(file).Decode(&data)
		if checkError(err) {
			return
		}
	}

	if data[folderID] != (time.Time{}) {
		lastDate = data[folderID]
	}
	return
}
