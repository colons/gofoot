package main

import "github.com/thoj/go-ircevent"

type Command interface {
	Initialize()
	Handle(*irc.Connection, *irc.Event)
	ShouldHandle(*irc.Event) bool
}

func getTarget(e *irc.Event) (target string) {
	if len(e.Arguments) > 0 {
		if e.Arguments[0] == nick {
			target = e.Nick
		} else {
			target = e.Arguments[0]
		}
	}
	return
}
