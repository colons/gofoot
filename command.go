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

func getTarget(e *irc.Event) (target string) {
	if len(e.Arguments) > 0 {
		if e.Arguments[0] == Con.GetNick() {
			target = e.Nick
		} else {
			target = e.Arguments[0]
		}
	}
	return
}

// Send a list of strings to target in a nicely formatted way
func sendStuff(target string, stuff []string) {
	Con.Privmsg(target, strings.Join(stuff, " \x034|\x03 "))
}

// take a message and split it into as many chunks as we expect args
func splitArgs(c ArgCommand, message string) []string {
	return strings.SplitN(message, " ", len(c.Args))
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

func argCommandShouldHandle(c ArgCommand, e *irc.Event) (bool) {
	if !strings.HasPrefix(e.Message, "!") {
		return false
	}

	userCommand := strings.TrimPrefix(e.Message, "!")

	split := splitArgs(c, userCommand)
	interesting := c.indecesOfSyntaxArguments()

	for _, interestingIndex := range interesting {
		expected := c.Args[interestingIndex]
		if split[interestingIndex] != expected {
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
