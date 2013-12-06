package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
)

type KatyCommand struct {
	ArgCommand
	KatyArray [6]string
}

func Katy() KatyCommand {
	return KatyCommand{
		ArgCommand{
			Args: []string{"katy", "[otaku]"},
			docs: "Expresses affection. Verbose.",
		},
		[6]string{
			"hi every1 im new!!!!!!! *holds up spork* my name is katy but u can call me t3h PeNgU1N oF d00m!!!!!!!! lol...as u can see im very random!!!! thats why i came here, 2 meet random ppl like me ^_^...",
			"im 13 years old (im mature 4 my age tho!!) i like 2 watch  invader zim w/ my girlfreind (im bi if u dont like it deal w/it) its our favorite tv show!!! bcuz its SOOOO random!!!!",
			"shes random 2 of course but i want 2 meet more random ppl =) like they say the more the merrier!!!! lol...neways i hope 2 make alot of freinds here so give me lots of commentses!!!!",
			"DOOOOOMMMM!!!!!!!!!!!!!!!! <--- me bein random again ^_^ hehe...toodles!!!!!",
			"love and waffles,",
			"* ~t3h PeNgU1N oF d00m~*",
		},
	}
}

func (c KatyCommand) Handle(e *irc.Event) {
	otaku := c.argsForCommand(e.Message)["otaku"]
	for _, k := range c.KatyArray {
		Connection.Privmsg(getTarget(e), strings.Replace(k, "katy", otaku, -1))
	}
}
