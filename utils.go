package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
)

// get the place we should be sending our response to event to
func getTarget(event *irc.Event) (target string) {
	if len(event.Arguments) > 0 {
		if event.Arguments[0] == Con.GetNick() {
			target = event.Nick
		} else {
			target = event.Arguments[0]
		}
	}
	return
}

// Send stuff to target in a nicely formatted way
func sendNestedStuff(target string, stuff [][]string) {
	subStuff := []string{}
	for _, thing := range(stuff) {
		subStuff = append(subStuff, strings.Join(thing, " \x034:\x03 "))
	}
	Con.Privmsg(target, strings.Join(subStuff, " \x034|\x03 "))
}

// Send stuff to target in a nicely formatted way
func sendStuff(target string, stuff []string) {
	Con.Privmsg(target, strings.Join(stuff, " \x034|\x03 "))
}
