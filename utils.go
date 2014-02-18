package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
)

// get the place we should be sending our response to event to
func getTarget(event *irc.Event) (target string) {
	if len(event.Arguments) > 0 {
		if event.Arguments[0] == Connection.GetNick() {
			target = event.Nick
		} else {
			target = event.Arguments[0]
		}
	}
	return
}

// strip consecutive whitespace from a string
func safeStuff(unsafe string) string {
	return strings.Join(strings.Fields(unsafe), " ")
}

func prettyNestedStuff(stuff [][]string) string {
	subStuff := []string{}
	for _, thing := range(stuff) {
		subStuff = append(subStuff, strings.Join(thing, " \x034:\x03 "))
	}
	return safeStuff(strings.Join(subStuff, " \x034|\x03 "))
}

func prettyStuff(stuff []string) string {
	return safeStuff(strings.Join(stuff, " \x034|\x03 "))
}
