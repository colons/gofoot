package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
)


type HelpQueryCommand struct {
	ArgCommand;
}

func HelpQuery() HelpQueryCommand {
	return HelpQueryCommand{
		ArgCommand{
			Args: []string{"help", "[query]"},
			docs: "Get information about a particular command.",
		},
	}
}

func (c HelpQueryCommand) Handle(e *irc.Event) {
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
