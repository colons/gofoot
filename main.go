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

