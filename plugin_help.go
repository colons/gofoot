package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
	"net/url"
)


type HelpCommand struct {
	ArgCommand
}

func Help() HelpCommand {
	return HelpCommand{
		ArgCommand{
			Args: []string{"help", "[query]"},
			docs: "Gets information about a particular command.",
		},
	}
}

func (c HelpCommand) Handle(e *irc.Event) {
	query := c.argsForCommand(e.Message)["query"]
	relaventCommands := []string{}

	for _, command := range(argCommands) {
		argCommand := argCommandFor(command)

		if argCommand.ShouldHandleMessage(e, query, false) {
			humanReadableArgs := Config.Event(e, "comchar") + strings.Join(argCommand.Args, " ")
			relaventCommands = append(relaventCommands, humanReadableArgs)
		}
	}

	if len(relaventCommands) > 0 {
		Connection.Privmsg(getTarget(e), prettyStuff(relaventCommands))
	}
}


type AboutCommand struct{
	ArgCommand
}

func About() AboutCommand {
	return AboutCommand{
		ArgCommand{
			Args: []string{"help"},
			docs: "Some useful links for your consideration.",
		},
	}
}

func (c AboutCommand) Handle(e *irc.Event) {
	helpUrlComponents := []string{
		Config.Event(e, "url"),
		"help",
		url.QueryEscape(Config.ourNetwork),
		url.QueryEscape(getTarget(e)),
	}

	Connection.Privmsg(getTarget(e), prettyNestedStuff([][]string{
		[]string{"\x02features\x02", strings.Join(helpUrlComponents, "/")},
		[]string{"\x02github\x02", "https://github.com/colons/gofoot"},
	}))
}
