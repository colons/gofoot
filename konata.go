package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
)


type KonataCommand struct {
	ArgCommand;
	KonataArray [4]string;
}

func Konata() KonataCommand {
	instance := KonataCommand{
		ArgCommand{Args: []string{"konata", "[nick]"}},
		[4]string{
			"I like konata because she is a otaku like me, except she has friends. Oh god I wish I had friends too ;_;",
			"konata also likes videogames and she is kawaii. And there are lesbians in the show and that's good because I like lesbians and I will never have a girlfriend. Why am I such a loser?!",
			"konata is like my dreamgirl she has a :3 face I love that. She is also nice why aren't real girls nice!? I got dumped a lot of times but I love konata and she wouldn't dump me because she's so  nice and cool.",
			"We would play videogames all day and watch Naruto and other cool animes on TV, and I would have sex with her because sex is so good. I wish I could have sex with a girl",
		},
	}
	return instance
}

func (c KonataCommand) ShouldHandle(e *irc.Event) bool {
	return argCommandShouldHandle(c.ArgCommand, e)
}

func (c KonataCommand) Handle(e *irc.Event) {
	nick := c.argsForCommand(e.Message)["nick"]
	for _, k := range(c.KonataArray) {
		Con.Privmsg(getTarget(e), strings.Replace(k, "konata", nick, -1))
	}
}