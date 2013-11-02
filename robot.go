package main

import (
	"strings"
	"fmt"
	"github.com/thoj/go-ircevent"
)

func RunRobot(network string) {
	Config = GetConfig(network)
	InitPersist(network)
	
	if Config.Network("address") == "" {
		fmt.Println("No address configured.")
		return
	}

	Connection = irc.IRC(Config.Network("nick"), Config.Network("user"))

	Connection.AddCallback("001", func(e *irc.Event) {
		nickservPassword := Config.Network("nickserv_password")
		if nickservPassword != "" {
			Connection.Privmsg("nickserv", "identify " + nickservPassword)
		}
		joinRooms()
	})

	Connection.AddCallback("NOTICE", func(e *irc.Event) {
		if strings.ToLower(e.Nick) == "nickserv" {
			joinRooms();
		}
	})

	Connection.AddCallback("INVITE", func(e *irc.Event) {
		Connection.Join(e.Message)
	})

	Connection.ReplaceCallback("CTCP_VERSION", 0, func(e *irc.Event) {
		Connection.SendRawf("NOTICE %s :\x01VERSION %s\x01", e.Nick, VERSION)
	})

	Connection.AddCallback("PRIVMSG", handleArgCommands)
	Connection.AddCallback("PRIVMSG", handleUnmanagedCommands)

	err := Connection.Connect(Config.Network("address"))
	if err != nil {
		fmt.Printf("Failed to connect to %s: %s", Config.Network("address"), err)
		return
	}

	Connection.Loop()
}

func joinRooms() {
	rooms := strings.Split(Config.Network("rooms"), ",")
	for _, room := range(rooms) {
		Connection.Join(room)
	}
}
