package main

import (
	"github.com/thoj/go-ircevent"
	"regexp"
)

var woofMatch *regexp.Regexp

type Woof struct {}

func (p Woof) Initialize() {
	woofMatch = regexp.MustCompile(".*\\b(woof)\\b.*")
}

func (p Woof) ShouldHandle(e *irc.Event) bool {
	return woofMatch.MatchString(e.Message)
}

func (p Woof) Handle(e *irc.Event) {
	Con.Privmsg(getTarget(e), "woof")
}
