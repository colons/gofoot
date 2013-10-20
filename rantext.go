package main

import (
	"github.com/thoj/go-ircevent"
	"os"
	"fmt"
	"bufio"
	"math/rand"
)

type Corpus struct {
	Text []string;
	Wrapper string;
}

type RantextCommand struct {
	ArgCommand;
	*Corpus;
}

func Rantext() (rantexts []CommandInterface) {
	sources := map[string]string{
		"jerkcity": "\x02%s\x02",
		"troll": "\x0304,08%s\x03",
	}

	for source, wrapper := range(sources) {
		corpus := getCorpus(source, wrapper)
		rantexts = append(rantexts, DirectedRantext(source, corpus))
		rantexts = append(rantexts, UndirectedRantext(source, corpus))
	}
	return
}

func DirectedRantext(source string, corpus *Corpus) RantextCommand {
	instance := RantextCommand{
		ArgCommand: ArgCommand{Args: []string{source, "[subject]"}},
		Corpus: corpus,
	}

	return instance
}


func UndirectedRantext(source string, corpus *Corpus) RantextCommand {
	instance := RantextCommand{
		ArgCommand: ArgCommand{Args: []string{source}},
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
	item := corpus.Text[rand.Intn(len(corpus.Text))]
	return fmt.Sprintf(corpus.Wrapper, item)
}


func (c RantextCommand) Handle(e *irc.Event) {
	subject := c.argsForCommand(e.Message)["subject"]
	var message string

	if subject != "" {
		message = fmt.Sprintf("%s: %s", subject, c.Corpus.Choice())
	} else {
		message = fmt.Sprintf("%s", c.Corpus.Choice())
	}

	Connection.Privmsg(getTarget(e), message)
}
