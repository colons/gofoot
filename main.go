package main

import (
    "github.com/thoj/go-ircevent"
    "fmt"
		"strings"
		"reflect"
		"math/rand"
		"time"
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
var argCommands []CommandInterface
var unmanagedCommands []CommandInterface

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	TheConfig = GetConfig()
	Con = irc.IRC(TheConfig.Nick, TheConfig.User)

	commands := []CommandInterface{
		Woof(), Http(), Konata(),
	}

	commands = append(commands, Rantext()...)

	argCommands = []CommandInterface{}
	unmanagedCommands = []CommandInterface{}

	for _, command := range(commands) {
		if (reflect.ValueOf(command).FieldByName("ArgCommand") == reflect.Value{}) {
			unmanagedCommands = append(unmanagedCommands, command)
		} else {
			argCommands = append(argCommands, command)
		}
	}

	Con.AddCallback("001", func(e *irc.Event) {
		for _, room := range(TheConfig.Rooms) {
			Con.Join(room)
		}
	})

	Con.ReplaceCallback("CTCP_VERSION", 0, func(e *irc.Event) {
		Con.SendRawf("NOTICE %s :\x01VERSION %s\x01", e.Nick, VERSION)
	})

	Con.AddCallback("PRIVMSG", handleArgCommands)
	Con.AddCallback("PRIVMSG", handleUnmanagedCommands)


	err := Con.Connect(TheConfig.Addr)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %s", TheConfig.Addr, err)
		return
	}

	Con.Loop()
}


// Return a Command's ArgCommand field after asserting that it is an ArgCommand
func argCommandFor(command CommandInterface) ArgCommand {
	return reflect.ValueOf(command).FieldByName("ArgCommand").Interface().(ArgCommand)
}

func handleArgCommands(e *irc.Event) {
	handlers := []CommandInterface{}

	for _, command := range(argCommands) {
		if command.ShouldHandle(e) {
			handlers = append(handlers, command)
		}
	}

	switch matches := len(handlers); matches {
	case 1:
		handlers[0].Handle(e)
	case 0:
		// no handlers :<
	default:
		disambiguator := []string{}
		for _, command := range(handlers) {
			args := argCommandFor(command).Args
			disambiguator = append(disambiguator, "!" + strings.Join(args, " "))
		}
		sendNestedStuff(getTarget(e), [][]string{[]string{"\x02ambiguous command\x02"}, disambiguator})
	}
}


func handleUnmanagedCommands(e *irc.Event) {
	for _, command := range(unmanagedCommands) {
		if command.ShouldHandle(e) {
			command.Handle(e)
		}
	}
}
