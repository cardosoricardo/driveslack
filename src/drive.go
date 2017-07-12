package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var relationMap map[string]string

var (
	typeFolder = "application/vnd.google-apps.folder"
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
	Version      string `json:"version"`
	OwnerUpdated string `json:"lastModifyingUserName"`
	URL          string `json:"alternateLink"`
	Type         string `json:"mimeType"`
}

//GetResponseFolder obtains the information about folder of google drive
func GetResponseFolder(folderID string) {
	//fmt.Println("entro")
	resp, err := http.Get("https://www.googleapis.com/drive/v2/files?q=%27" + folderID + "%27+in+parents&key=" + tokenDrive)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var drive DriveFolder

	json.Unmarshal(body, &drive)

	//var countFolders int
	for _, file := range drive.Items {
		if file.Type == typeFolder {
			GetResponseFolder(file.ID)
			continue
		}
		fmt.Println("file")
	}
	//fmt.Println("Hay " + strconv.Itoa(countFolders) + " folders y " + strconv.Itoa((len(drive.Items) - countFolders)) + " archivos.")
}
