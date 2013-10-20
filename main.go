package main

import (
	"github.com/thoj/go-ircevent"
	"fmt"
	"strings"
	"reflect"
	"math/rand"
	"time"
	"os"
)

const (
	VERSION = "gofoot"
	USAGE = "gofoot robot <network>"
)

var Connection *irc.Connection
var argCommands []CommandInterface
var unmanagedCommands []CommandInterface
var Config config

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println(USAGE)
		return
	}

	switch mode := os.Args[1]; mode {
	case "robot":
		if len(os.Args) < 3 {
			fmt.Println(USAGE)
			return
		} else {
			network := os.Args[2]
			RunRobot(network)
		}
	// XXX case "server":
	}
}

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
