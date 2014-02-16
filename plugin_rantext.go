package main

import (
	"github.com/thoj/go-ircevent"
	"os"
	"fmt"
	"bufio"
	"math/rand"
	"strings"
)

type Corpus struct {
	Text []string
	Wrapper string
}

type RantextCommand struct {
	ArgCommand
	*Corpus
}

func Rantext() (rantexts []CommandInterface) {
	wrappers := map[string]string{
		"jerkcity": "\x02%s",
		"troll": "\x0304,08%s",
		"sharon": "\x0302,00%s",
		"thatcher": "\x0302%s",
		"theo": "\x0304,01%s",
		"toothbrush": "\x02\x0307%s",
	}

	sources := strings.Split(Config.Global("rantext_sources"), ",")

	for _, source := range(sources) {
		wrapper := wrappers[source]
		if wrapper == "" {
			wrapper = "%s"
		}
		corpus := getCorpus(source, wrapper)
		rantexts = append(rantexts, DirectedRantext(source, corpus))
		rantexts = append(rantexts, UndirectedRantext(source, corpus))
	}
	return
}

func DirectedRantext(source string, corpus *Corpus) RantextCommand {
	instance := RantextCommand{
		ArgCommand: ArgCommand{
			Args: []string{source, "[subject]"},
		},
		Corpus: corpus,
	}

	return instance
}

func (r RantextCommand) GetDocs() string {
	var theFmt string
	if len(r.Args) == 2 {
		theFmt = "`[subject]: %s`"
	} else {
		theFmt = "`%s`"
	}
	return fmt.Sprintf(theFmt, r.Corpus.Choice())
}


func UndirectedRantext(source string, corpus *Corpus) RantextCommand {
	instance := RantextCommand{
		ArgCommand: ArgCommand{
			Args: []string{source},
		},
		Corpus: corpus,
	}

	return instance
}


func getCorpus(source, wrapper string) (*Corpus) {
	pathFormat := os.Getenv("HOME") + "/.gofoot/rantext/%s.txt"
	file, err := os.Open(fmt.Sprintf(pathFormat, source))
	if err != nil {
		fmt.Println("Could not open %s: %s", source, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	corpus := Corpus{}

	for scanner.Scan() {
		corpus.Text = append(corpus.Text, scanner.Text())
	}

	corpus.Wrapper = wrapper
	return &corpus
}


func (corpus Corpus) Choice() string {
	return strings.TrimSpace(corpus.Text[rand.Intn(len(corpus.Text))])
}

func (corpus Corpus) WrappedChoice() string {
	return fmt.Sprintf(corpus.Wrapper, corpus.Choice())
}


func (c RantextCommand) Handle(e *irc.Event) {
	subject := c.argsForCommand(e.Message)["subject"]
	choice := c.Corpus.WrappedChoice()
	var message string

	if subject != "" {
		message = fmt.Sprintf("%s: %s", subject, choice)
	} else {
		message = fmt.Sprintf("%s", choice)
	}

	Connection.Privmsg(getTarget(e), message)
}
