package main

import (
	"flag"
	"io/ioutil"
	"os"
	"time"

	"fmt"

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
	minutes         *int
)

func init() {
	dataConfig := flag.String("c", "", "path to YAML configuration")
	relation = flag.String("d", "", "path to relation driveID and channelID")
	minutes = flag.Int("t", 15, "time to execute the periodic task in minutes")
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
	periodicTime := *minutes
	ticker := time.NewTicker(time.Duration(periodicTime) * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				relationArray := GetRelationFromFile(*relation)
				if len(relationArray) == 0 {
					panic("relations not register")
				}

				for _, relation := range relationArray {
					fmt.Println("request")
					lastDate := GetResponseFolder(relation.DriveID, relation.ChannelID, time.Time{}, true)
					if lastDate == (time.Time{}) {
						return
					}
					Save(relation.DriveID, lastDate)

				}

				ticker = time.NewTicker(time.Duration(periodicTime) * time.Minute)
			}
		}
	}()

	quit := make(chan bool, 1)
	// main will continue to wait untill there is an entry in quit
	<-quit
}
