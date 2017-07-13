package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

var relationMap map[string]string

var (
	typeFolder = "application/vnd.google-apps.folder"
)

var (
	templateNewFolder    = "%s creó la carpeta %s en %s"
	templateNewFile      = "%s creó el archivo %s  en %s"
	templateUpdateFile   = "%s actualizó el archivo %s en %s"
	templateUpdateFolder = "%s actualizó el folder %s en %s"
)

//DriveFolder define a struct of drive folder response
type DriveFolder struct {
	Items []ItemDrive `json:"items"`
}

//ItemDrive contains information about file
type ItemDrive struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Updated      time.Time `json:"modifiedDate"`
	Created      time.Time `json:"createdDate"`
	OwnerUpdated string    `json:"lastModifyingUserName"`
	URL          string    `json:"alternateLink"`
	Type         string    `json:"mimeType"`
}

//GetResponseFolder obtains the information about folder of google drive
func GetResponseFolder(folderID, channelID string, lastUpdated time.Time, root bool) time.Time {

	var drive DriveFolder
	//var lastDate string

	resp, err := http.Get("https://www.googleapis.com/drive/v2/files?q=%27" + folderID + "%27+in+parents&key=" + tokenDrive)
	if checkError(err) {
		return time.Time{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if checkError(err) {
		return time.Time{}
	}
	json.Unmarshal(body, &drive)

	if len(drive.Items) == 0 {
		return time.Time{}
	}
	if lastUpdated == (time.Time{}) { //last date saves in db by folderID
		lastUpdated = Get(folderID)
	}
	maxDate := lastUpdated
	for _, file := range drive.Items {
		template := getTemplate(file)
		dateFile := file.Updated
		if dateFile.After(lastUpdated) {
			ex := RegisterMessage(template, channelID, file.OwnerUpdated, file.Title, file.URL)
			if checkError(ex) {
				return time.Time{}
			}
		}
		if dateFile.After(maxDate) {
			maxDate = dateFile
		}

		if file.Type == typeFolder {
			lastDate := GetResponseFolder(file.ID, channelID, lastUpdated, false)
			if lastDate.After(maxDate) {
				maxDate = lastDate
			}
		}
	}
	return maxDate
}

//getTemplate obatins template message of slack
func getTemplate(file ItemDrive) (template string) {

	if file.Created == file.Updated {
		switch file.Type {
		case typeFolder:
			template = templateNewFolder
			break
		default:
			template = templateNewFile
			break
		}
	} else {
		switch file.Type {
		case typeFolder:
			template = templateUpdateFolder
			break
		default:
			template = templateUpdateFile
			break
		}
	}

	return
}
