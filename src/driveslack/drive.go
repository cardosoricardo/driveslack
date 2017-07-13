package main

import (
	"encoding/json"
	"fmt"
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
	ID           string `json:"id"`
	Title        string `json:"title"`
	Updated      string `json:"modifiedDate"`
	Created      string `json:"createdDate"`
	Version      string `json:"version"`
	OwnerUpdated string `json:"lastModifyingUserName"`
	URL          string `json:"alternateLink"`
	Type         string `json:"mimeType"`
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
	//lastDate = drive.Items[0].Updated
	maxDate := lastUpdated
	//fmt.Println(lastUpdated)
	for _, file := range drive.Items {
		//template := getTemplate(file)
		dateFile, _ := time.Parse(time.RFC3339, file.Updated)
		if dateFile.After(lastUpdated) {
			fmt.Println("message", file.Title)
			// RegisterMessage(template, channelID, file.OwnerUpdated, file.Title, file.URL)
		}
		if dateFile.After(maxDate) {
			maxDate = dateFile
		}

		// if (dateFile.Before(lastUpdated) || dateFile.Equal(lastUpdated)) && file.Type != typeFolder { //compare two times
		// 	break
		// }
		//fmt.Println("message")

		if file.Type == typeFolder {
			lastDate := GetResponseFolder(file.ID, channelID, lastUpdated, false)
			if lastDate.After(maxDate) {
				maxDate = lastDate
			}
		}
	}
	// fmt.Println(maxDate)
	if root {
		f, err := time.Parse(time.UnixDate, maxDate.String())
		fmt.Println(err)
		Save(folderID, f.Format(time.RFC3339))
	}
	return maxDate
}

//getTemplate obatins template message of slack
func getTemplate(file ItemDrive) (template string) {

	t1, _ := time.Parse(time.RFC3339, file.Created)
	t2, _ := time.Parse(time.RFC3339, file.Updated)
	if t1 == t2 {
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
