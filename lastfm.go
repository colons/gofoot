package main

import (
	"github.com/thoj/go-ircevent"
	"encoding/json"
	"net/http"
	"fmt"
	"net/url"
	"io"
)

type NowPlayingForUserCommand struct{
	ArgCommand
}

type RecentTrack struct {
	Artist map[string]string
	Name string
	Album map[string]string
	Attr map[string]string `json:"@attr"`
}

type NPResponse struct{
	RecentTracks map[string][]RecentTrack
}

func NowPlayingForUser() NowPlayingForUserCommand {
	return NowPlayingForUserCommand{
		ArgCommand{
			Args: []string{"np", "[user]"},
			docs: "Show everyone what you're playing right now.",
		},
	}
}

func (c NowPlayingForUserCommand) Handle(e *irc.Event) {
	recentest := recentestForUser("NiviJh")
	trackData := []string{
		fmt.Sprintf("\x02%s\x02", recentest.Name),
	}

	album := recentest.Album["#text"]
	if album != "" {
		trackData = append(trackData, album)
	}

	artist := recentest.Artist["#text"]
	if artist != "" {
		trackData = append(trackData, artist)
	}

	userData := []string{fmt.Sprintf("http://last.fm/user/%s", "NiviJh")}

	data := [][]string{trackData, userData}
	Connection.Privmsg(getTarget(e), prettyNestedStuff(data))
}

func recentestForUser(user string) RecentTrack {
	urlFmt := "http://ws.audioscrobbler.com/2.0/?%s"

	args := url.Values{}
	args.Set("format", "json")
	args.Set("api_key", Config.Network("lastfm_api_key"))
	args.Set("method", "user.getRecentTracks")
	args.Set("user", "NiviJh")

	theUrl := fmt.Sprintf(urlFmt, args.Encode())
	fmt.Println(theUrl)

	resp, err := http.Get(theUrl)
	if err != nil {
		fmt.Printf("Could not np: %s", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var np NPResponse

	for {
		if err := dec.Decode(&np); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
	}
	
	return np.RecentTracks["track"][0]
}
