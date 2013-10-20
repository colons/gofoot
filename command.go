package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
)

type CommandInterface interface {
	Handle(*irc.Event)
	ShouldHandle(*irc.Event) bool
}

type ArgCommand struct {
	Args []string;
}

// take a message and split it into as many chunks as we expect args
func splitArgs(c ArgCommand, message string) []string {
	strippedMessage := strings.TrimRight(message, " ")
	return strings.SplitN(strippedMessage, " ", len(c.Args))
}

func argIsVariable(arg string) bool {
	return strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]")
}

// return the indeces of arguments we want to hand to handlers;
// for instance, a command with args
// ["translate", "[from]", "[to]", "[query]"] would return [1, 2, 3]
func (c ArgCommand) indecesOfVariableArguments() (interesting []int) {
	for i, arg := range c.Args {
		if argIsVariable(arg) {
			interesting = append(interesting, i)
		}
	}
	return
}

// Return the inverse of indecesOfVariableArguments; that is,
// [0] for the above example.
func (c ArgCommand) indecesOfSyntaxArguments() (interesting []int) {
	for i, arg := range c.Args {
		if !argIsVariable(arg) {
			interesting = append(interesting, i)
		}
	}
	return
}

func (c ArgCommand) ShouldHandle(e *irc.Event) bool {
	if !strings.HasPrefix(e.Message, "!") {
		return false
	}

	return c.ShouldHandleMessage(e.Message, true)
}

func (c ArgCommand) ShouldHandleMessage(message string, requireAllArguments bool) bool {
	userCommand := strings.TrimPrefix(message, "!")

	split := splitArgs(c, userCommand)
	interesting := c.indecesOfSyntaxArguments()

	if requireAllArguments {
		if len(split) != len(c.Args) {
			return false
		}
	}

	for _, i := range interesting {
		expected := c.Args[i]
		if !strings.HasPrefix(expected, split[i]) {
			return false
		}
	}

	return true
}


// Return a map of variable arguments to their value in a given use.
// In our example, return {"from": "fr", "to": "en", "query": "bonbon hahah"}
// given input "!tra fr en bonbonb hahah"
func (c ArgCommand) argsForCommand(command string) (args map[string]string) {
	variable := c.indecesOfVariableArguments()
	split := splitArgs(c, command)
	args = make(map[string]string)

	for _, i := range variable {
		name := strings.TrimPrefix(c.Args[i], "[")
		name = strings.TrimSuffix(name, "]")
		args[name] = split[i]
	}

	return
}
