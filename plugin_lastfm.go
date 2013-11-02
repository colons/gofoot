package main

import (
	"github.com/thoj/go-ircevent"
	"encoding/json"
	"net/http"
	"fmt"
	"net/url"
	"io"
)


const urlFmt = "http://ws.audioscrobbler.com/2.0/?%s"

type NowPlayingForUserCommand struct{
	ArgCommand
}

func NowPlayingForUser() NowPlayingForUserCommand {
	return NowPlayingForUserCommand{
		ArgCommand{
			Args: []string{"np", "[user]"},
			docs: "Spams the channel with someone else's terrible taste in music.",
		},
	}
}


type NowPlayingCommand struct{
	ArgCommand
}

func NowPlaying() NowPlayingCommand {
	return NowPlayingCommand{
		ArgCommand{
			Args: []string{"np"},
			docs: "Spams the channel with your terrible taste in music.",
		},
	}
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

func (c NowPlayingForUserCommand) Handle(e *irc.Event) {
	user := c.argsForCommand(e.Message)["user"]
	BroadcastNowPlayingFor(e, user)
}

func (c NowPlayingCommand) Handle(e *irc.Event) {
	user := Remember(fmt.Sprintf("lastfm_user:%s", e.Nick))
	if user == "" {
		Connection.Privmsg(getTarget(e), fmt.Sprintf("i dunno who you are :<"))
		return
	}
	BroadcastNowPlayingFor(e, user)
}

func BroadcastNowPlayingFor(e *irc.Event, user string) {
	recentest := recentestTrackForUser(user)
	if recentest.Name == "" {
		// this is probably not a real user
		Connection.Privmsg(getTarget(e), "nope")
		return
	}

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

	userData := []string{fmt.Sprintf("http://last.fm/user/%s", user)}

	data := [][]string{trackData, userData}
	Connection.Privmsg(getTarget(e), prettyNestedStuff(data))
}

func recentestTrackForUser(user string) RecentTrack {
	args := url.Values{}
	args.Set("format", "json")
	args.Set("api_key", Config.Network("lastfm_api_key"))
	args.Set("method", "user.getRecentTracks")
	args.Set("user", user)

	theUrl := fmt.Sprintf(urlFmt, args.Encode())

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
	
	if len(np.RecentTracks["track"]) > 0 {
		return np.RecentTracks["track"][0]
	}	else {
		return RecentTrack{}
	}
}


type LastFMUserInfo struct {
	Error int
}

func lastFMUserExists(user string) bool {
	args := url.Values{}
	args.Set("format", "json")
	args.Set("api_key", Config.Network("lastfm_api_key"))
	args.Set("method", "user.getInfo")
	args.Set("user", user)

	theUrl := fmt.Sprintf(urlFmt, args.Encode())

	resp, err := http.Get(theUrl)
	if err != nil {
		fmt.Printf("Could not get user: %s", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var userInfo LastFMUserInfo

	for {
		if err := dec.Decode(&userInfo); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
	}

	return userInfo.Error == 0
}


type SetLastFMCommand struct{
	ArgCommand
}


func SetLastFM() SetLastFMCommand {
	return SetLastFMCommand{
		ArgCommand{
			Args: []string{"set", "lastfm", "[user]"},
			docs: "Remembers who you are.",
		},
	}
}

func (c SetLastFMCommand) Handle(e *irc.Event) {
	username := c.argsForCommand(e.Message)["user"]
	if lastFMUserExists(username) {
		Persist(fmt.Sprintf("lastfm_user:%s", e.Nick), username)
		Connection.Privmsg(getTarget(e), fmt.Sprintf("\x02%s\x02 is now last.fm user \x02%s\x02", e.Nick, username))
	} else {
		Connection.Privmsg(getTarget(e), fmt.Sprintf("\x02%s\x02 does not exist", username))
	}
}
