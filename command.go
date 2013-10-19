package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
)

type Command interface {
	Initialize()
	Handle(*irc.Event)
	ShouldHandle(*irc.Event) bool
}

func getTarget(e *irc.Event) (target string) {
	if len(e.Arguments) > 0 {
		if e.Arguments[0] == Con.GetNick() {
			target = e.Nick
		} else {
			target = e.Arguments[0]
		}
	}
	return
}

// Send a list of strings to target in a nicely formatted way
func sendStuff(target string, stuff []string) {
	Con.Privmsg(target, strings.Join(stuff, " \x034|\x03 "))
}
