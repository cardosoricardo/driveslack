package main

import (
	"flag"
	"io/ioutil"
	"os"
	"time"

	"github.com/cardosoricardo/driveslack/src/driveslack/app"

	yaml "gopkg.in/yaml.v2"
)

var (
	conf            = app.Config{}
	relation, files *string
	minutes         *int
)

func init() {
	dataConfig := flag.String("c", "", "path to YAML configuration")
	relation = flag.String("r", "", "path to relation driveID and channelID JSON file")
	minutes = flag.Int("t", 15, "time to execute the periodic task in minutes")
	flag.Parse()
	if len(*dataConfig)*len(*relation) == 0 {
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

	app.SetVars(conf)
}

func main() {
	periodicTime := *minutes
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				relationArray, _ := app.GetRelationFromFile(*relation)
				for _, relation := range relationArray {
					lastDate, _ := app.GetResponseFolder(relation.DriveID, relation.ChannelID, time.Time{})
					if lastDate == (time.Time{}) {
						return
					}
					app.Save(relation.DriveID, lastDate)
				}

				ticker = time.NewTicker(time.Duration(periodicTime) * time.Minute)
			}
		}
	}()

	quit := make(chan bool, 1)
	<-quit
}
