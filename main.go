package main

import (
	"github.com/thoj/go-ircevent"
	"fmt"
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
var Commands []CommandInterface
var GlobalConfig config
var Config config

func main() {
	GlobalConfig = GetConfig("")
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println(USAGE)
		return
	}

	// initialize commands
	Commands = []CommandInterface{
		HelpQuery(), Woof(), Http(), Konata(),
	}
	Commands = append(Commands, Rantext()...)

	// divvy up commands into managed and unmanaged
	for _, command := range(Commands) {
		if IsArgCommand(command) {
			argCommands = append(argCommands, command)
		} else {
			unmanagedCommands = append(unmanagedCommands, command)
		}
	}

	// and do what we came here for
	switch mode := os.Args[1]; mode {
	case "robot":
		if len(os.Args) < 3 {
			fmt.Println(USAGE)
			return
		} else {
			network := os.Args[2]
			RunRobot(network)
		}
	case "server":
		RunServer()
	}
}
