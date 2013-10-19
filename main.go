package main

import (
    "github.com/thoj/go-ircevent"
    "fmt"
)

const (
	VERSION = "gofoot"
)

type Config struct {
	// define GetConfig to return one of these
	Rooms []string
	Nick string
	User string
	Addr string
}

var TheConfig Config
var Con *irc.Connection

func main() {
	TheConfig = GetConfig()
	Con = irc.IRC(TheConfig.Nick, TheConfig.User)
	commands := []Command{
		WoofCommand{}, HttpCommand{},
	}

	Con.AddCallback("001", func(e *irc.Event) {
		for i := 0; i < len(TheConfig.Rooms); i++ {
			Con.Join(TheConfig.Rooms[i])
		}
	})

	Con.ReplaceCallback("CTCP_VERSION", 0, func(e *irc.Event) {
		Con.SendRawf("NOTICE %s :\x01VERSION %s\x01", e.Nick, VERSION)
	})


	for i := 0; i < len(commands); i++ {
		command := commands[i]
		command.Initialize()

		Con.AddCallback("PRIVMSG", func(e *irc.Event) {
			if command.ShouldHandle(e) {
				command.Handle(e)
			}
		})
	}

	err := Con.Connect(TheConfig.Addr)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %s", TheConfig.Addr, err)
		return
	}

	Con.Loop()
}
