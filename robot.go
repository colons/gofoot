package main

import (
	"strings"
	"reflect"
	"fmt"
	"github.com/thoj/go-ircevent"
)

func RunRobot(network string) {
	Config = GetConfig(network)

	if Config.Network("address") == "" {
		fmt.Println("No address configured.")
		return
	}

	Connection = irc.IRC(Config.Network("nick"), Config.Network("user"))

	commands := []CommandInterface{
		HelpQuery(), Woof(), Http(), Konata(),
	}

	commands = append(commands, Rantext()...)

	for _, command := range(commands) {
		if (reflect.ValueOf(command).FieldByName("ArgCommand") == reflect.Value{}) {
			unmanagedCommands = append(unmanagedCommands, command)
		} else {
			argCommands = append(argCommands, command)
		}
	}

	Connection.AddCallback("001", func(e *irc.Event) {
		rooms := strings.Split(Config.Network("rooms"), ",")

		for _, room := range(rooms) {
			Connection.Join(room)
		}
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
