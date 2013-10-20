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
		ArgCommand{Args: []string{"help", "[query]"}},
	}
}

func (c HelpQueryCommand) Handle(e *irc.Event) {
	query := c.argsForCommand(e.Message)["query"]
	relaventCommands := []string{}

	for _, command := range(argCommands) {
		argCommand := argCommandFor(command)

		if argCommand.ShouldHandleMessage(query, false) {
			humanReadableArgs := "!" + strings.Join(argCommand.Args, " ")
			relaventCommands = append(relaventCommands, humanReadableArgs)
		}
	}

	if len(relaventCommands) > 0 {
		sendStuff(getTarget(e), relaventCommands)
	}
}
