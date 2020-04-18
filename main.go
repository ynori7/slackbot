package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/ynori7/slackbot/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("you must specify the path to the config file"))
	}

	//Get the config
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var conf config.Config
	if err := conf.Parse(data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)

	api := slack.New(conf.SlackToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		fmt.Printf("Event received %+v\n", msg)
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			fmt.Println("Connected!")

			chans, err := rtm.GetChannels(true)
			if err != nil {
				fmt.Println("Error getting channels", err)
				os.Exit(1)
			}
			for _, c := range chans {
				conf.Channels[c.ID] = c
			}
			for _, a := range conf.Admins {
				_, _, conf.AdminChannels[a], err = api.OpenIMChannel(a)
				if err != nil {
					fmt.Println("Error opening channel to admin", err)
					os.Exit(1)
				}
			}
		case *slack.MessageEvent:
			//fmt.Println(ev.Text)
			if ev.SubType != "message_replied" {
				fmt.Printf("%+v\n", ev)
				rtm.SendMessage(rtm.NewOutgoingMessage("blah", conf.AdminChannels[conf.Admins[0]]))
			}
		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
		}
	}
}
