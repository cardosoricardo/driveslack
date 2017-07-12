package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

var (
	iconURL    string
	tokenSlack string
	tokenDrive string
	username   string
)

//SetVars set tokens.
func SetVars() {
	iconURL = conf.IconURL
	tokenSlack = conf.Slack
	tokenDrive = conf.Drive
	username = conf.Username
}

//Channel define a struct of channel's slack
type Channel struct {
	ID          string
	Description string
}

//RegisterMessage register a mmesage in channel
func RegisterMessage(template, channelID, username, filename, url string) {
	message := fmt.Sprintf(template, username, filename, url)
	api := slack.New(tokenSlack)
	params := slack.PostMessageParameters{
		IconURL:  iconURL,
		Username: username,
	}
	_, _, err := api.PostMessage(channelID, message, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf(message)
}

//GetChannels obtains the list of his channels
func GetChannels() {
	var channels []Channel
	api := slack.New(tokenSlack)
	response, err := api.GetChannels(true)
	if checkError(err) {
		return
	}
	for _, c := range response {
		channel := Channel{
			ID:          c.ID,
			Description: c.Name,
		}

		channels = append(channels, channel)
	}
	// pretty.Println(channels)
}

//GetGroups obtains the list of his private channels
func GetGroups() {
	var channels []Channel
	api := slack.New(tokenSlack)
	response, err := api.GetGroups(true)
	if checkError(err) {
		return
	}
	for _, c := range response {
		channel := Channel{
			ID:          c.ID,
			Description: c.Name,
		}

		channels = append(channels, channel)
	}
	// pretty.Println(channels)
}

func checkError(err error) bool {
	if err != nil {
		fmt.Printf("%s\n", err)
		return true
	}
	return false
}
