package main

import (
	"github.com/thoj/go-ircevent"
	"regexp"
)

var woofMatch *regexp.Regexp

type WoofCommand struct {}

func Woof() WoofCommand {
	woofMatch = regexp.MustCompile(".*\\b(woof)\\b.*")
	instance := WoofCommand{}
	return instance
}

func (c WoofCommand) ShouldHandle(e *irc.Event) bool {
	return woofMatch.MatchString(e.Message)
}

func (c WoofCommand) Handle(e *irc.Event) {
	Connection.Privmsg(getTarget(e), "woof")
}
