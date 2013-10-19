package main

import (
    "github.com/thoj/go-ircevent"
    "fmt"
)

const (
	VERSION = "gofoot"
)

var rooms = []string{"#test"}
var nick = "gofoot"
var addr = "irc.99chan.org:6667"

func main() {
	con := irc.IRC(nick, nick)
	err := con.Connect(addr)

	if err != nil {
		fmt.Printf("Failed to connect to %s: %s", addr, err)
		return
	}

	commands := []Command{Woof{}}

	for i := 0; i < len(commands); i++ {
		commands[i].Initialize()
	}

	con.AddCallback("001", func(e *irc.Event) {
		for i := 0; i < len(rooms); i++ {
			con.Join(rooms[i])
		}
	})

	con.ReplaceCallback("CTCP_VERSION", 0, func(e *irc.Event) {
		con.SendRawf("NOTICE %s :\x01VERSION %s\x01", e.Nick, VERSION)
	})

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		for i := 0; i < len(commands); i++ {
			command := commands[i]
			if command.ShouldHandle(e) {
				command.Handle(con, e)
			}
		}
	})

	con.Loop()
}
