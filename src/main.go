package main

import (
	"flag"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type data struct {
	Slack    string `yaml:"slack-token"`
	Drive    string `yaml:"drive-token"`
	IconURL  string `yaml:"icon-url"`
	Username string `yaml:"username"`
}

//DriveSlack defines a struct of relation between drive id and slack id
type DriveSlack struct {
	ChannelID string `json:"channel"`
	DriveID   string `json:"drive"`
}

var (
	conf = data{}
	info *string
)

func init() {
	dataConfig := flag.String("data", "", "path to YAML configuration")
	info = flag.String("info", "", "path to file information")
	flag.Parse()
	if len(*dataConfig)*len(*info) <= 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	bytes, err := ioutil.ReadFile(*dataConfig)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		panic(err)
	}
	SetVars()
}

func main() {
	message := "He adjuntado un archivo http://www.example.com.mx"
	RegisterMessage("random", message)

	// GetChannels()

	// GetGroups()

	// GetFilesByChannel(*info)
}
