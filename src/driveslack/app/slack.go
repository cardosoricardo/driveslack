package app

import (
	"fmt"

	"github.com/nlopes/slack"
)

var (
	conf        Config
	slackClient *slack.Client
	slackParams slack.PostMessageParameters
)

//Config defines a struct with configuration parameters of App.
type Config struct {
	Slack    string `yaml:"slack-token"`
	Drive    string `yaml:"drive-token"`
	IconURL  string `yaml:"icon-url"`
	Username string `yaml:"username"`
}

//SetVars set tokens.
func SetVars(c Config) {
	conf = c
	slackClient = slack.New(conf.Slack)
	slackParams = slack.PostMessageParameters{
		IconURL:  conf.IconURL,
		Username: conf.Username,
	}
}

//Channel define a struct of channel's slack
type Channel struct {
	ID          string
	Description string
}

//registerMessage register a mmesage in channel
func registerMessage(template, channelID, username, filename, url string) error {
	message := fmt.Sprintf(template, username, filename, url)
	_, _, err := slackClient.PostMessage(channelID, message, slackParams)
	checkError(err)
	return err
}

//getChannels obtains the list of his channels
func getChannels() (channels []Channel) {
	response, err := slackClient.GetChannels(true)
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
	return
}

//getGroups obtains the list of his private channels
func getGroups() (channels []Channel) {
	response, err := slackClient.GetGroups(true)
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
	return channels
}
