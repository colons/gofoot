package main

import (
	"github.com/thoj/go-ircevent"
	"regexp"
)

var woofMatch *regexp.Regexp

type WoofCommand struct {}

func (c WoofCommand) Initialize() {
	woofMatch = regexp.MustCompile(".*\\b(woof)\\b.*")
}

func (c WoofCommand) ShouldHandle(e *irc.Event) bool {
	return woofMatch.MatchString(e.Message)
}

func (c WoofCommand) Handle(e *irc.Event) {
	Con.Privmsg(getTarget(e), "woof")
}
