package main

import (
	"strings"
	"fmt"
	"github.com/thoj/go-ircevent"
	"os/signal"
	"os"
)

func RunRobot(network string) {
	Config = GetConfig(network)
	InitPersist()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			DB.Close()
			fmt.Println(sig)
			os.Exit(0)
		}
	}()

	
	if Config.Network("address") == "" {
		fmt.Println("No address configured.")
		return
	}

	Connection = irc.IRC(Config.Network("nick"), Config.Network("user"))

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
