package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
	"reflect"
	"fmt"
)

type CommandInterface interface {
	Handle(*irc.Event)
	ShouldHandle(*irc.Event) bool
	GetDocs() string
}


type ArgCommand struct {
	Args []string
	docs string
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
	if !strings.HasPrefix(e.Message, Config.Event(e, "comchar")) {
		return false
	}

	return c.ShouldHandleMessage(e, e.Message, true)
}


// Decide whether we should handle message. Takes an irc.Event for config's sake
func (c ArgCommand) ShouldHandleMessage(e *irc.Event, message string, requireAllArguments bool) bool {
	userCommand := strings.TrimPrefix(message, Config.Event(e, "comchar"))

	split := splitArgs(c, userCommand)
	interesting := c.indecesOfSyntaxArguments()

	if requireAllArguments {
		if len(split) != len(c.Args) {
			return false
		}
	}

	for _, i := range interesting {
		expected := c.Args[i]
		var part string
		if !requireAllArguments && i >= len(split) {
			part = ""
		} else {
			part = split[i]
		}

		if !strings.HasPrefix(expected, part) {
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

func (c ArgCommand) GetDocs() string {
	return c.docs
}


// Return a Command's ArgCommand field after asserting that it is an ArgCommand
func argCommandFor(command CommandInterface) ArgCommand {
	return reflect.ValueOf(command).FieldByName("ArgCommand").Interface().(ArgCommand)
}


// Decide if a command should be allowed to handle an event based on the configured blacklist
func commandIsAllowedToHandle(command CommandInterface, e *irc.Event) bool {
	commandName := fmt.Sprint(reflect.ValueOf(command).Type())
	commandName = strings.TrimPrefix(commandName, "main.")
	commandName = strings.TrimSuffix(commandName, "Command")

	for _, blacklisted := range strings.Split(Config.Event(e, "blacklist"), ",") {
		if commandName == blacklisted {
			return false
		}
	}
	for _, ignored := range strings.Split(Config.Event(e, "ignore"), ",") {
		if e.Nick == ignored {
			return false
		}
	}
	return true
}


func handleArgCommands(e *irc.Event) {
	handlers := []CommandInterface{}

	for _, command := range(argCommands) {
		if commandIsAllowedToHandle(command, e) && command.ShouldHandle(e) {
			handlers = append(handlers, command)
		}
	}

	switch matches := len(handlers); matches {
	case 1:
		handlers[0].Handle(e)
	case 0:
		// no handlers :<
	default:
		disambiguator := []string{}
		for _, command := range(handlers) {
			args := argCommandFor(command).Args
			disambiguator = append(disambiguator, "!" + strings.Join(args, " "))
		}

		pretty := prettyNestedStuff([][]string{[]string{"\x02ambiguous command\x02"}, disambiguator})
		Connection.Privmsg(getTarget(e), pretty)
	}
}


func handleUnmanagedCommands(e *irc.Event) {
	for _, command := range(unmanagedCommands) {
		if commandIsAllowedToHandle(command, e) && command.ShouldHandle(e) {
			command.Handle(e)
		}
	}
}

func IsArgCommand(command CommandInterface) bool {
	return reflect.ValueOf(command).FieldByName("ArgCommand") != reflect.Value{}
}
