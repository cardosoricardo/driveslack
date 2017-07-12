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
	conf            = data{}
	relation, files *string
)

func init() {
	dataConfig := flag.String("data", "", "path to YAML configuration")
	relation = flag.String("info", "", "path to relation driveID and channelID")
	flag.Parse()
	if len(*dataConfig)*len(*relation) <= 0 {
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
	// message := "He adjuntado un archivo http://www.example.com.mx"
	// RegisterMessage("random", message)

	// GetChannels()

	// GetGroups()

	relationArray := GetRelationFromFile(*relation)
	if len(relationArray) == 0 {
		panic("No hay relaciones registradas")
	}
	GetResponseFolder("0B9kaSejrgCGDYnNxZzRaWFhhUWs")
}
