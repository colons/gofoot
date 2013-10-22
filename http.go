package main

import (
	"github.com/thoj/go-ircevent"
	"regexp"
	"net/http"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go-html-transform/html/transform"
	"fmt"
	"strings"
)

var urlMatch *regexp.Regexp
var titleSelector *selector.Chain

type HttpCommand struct {}

func Http() HttpCommand {
	urlMatch = regexp.MustCompile("\\b(https?://\\S+)\\b")
	instance := HttpCommand{}
	return instance
}


func (c HttpCommand) ShouldHandle(e *irc.Event) bool {
	return (urlMatch.FindStringIndex(e.Message) != nil)
}


func (c HttpCommand) Handle(e *irc.Event) {
	urls := urlMatch.FindAllString(e.Message, -1)

	for _, url := range urls {
		metadata := metadataForUrl(url)
		if metadata != nil {
			sendStuff(getTarget(e), metadata)
		}
	}
}


func metadataForUrl(url string) []string {
	metadata := []string{}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting %s: %s\n", url, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		metadata = append(metadata, fmt.Sprintf("%d", resp.StatusCode))
	}

	contentType := strings.Split(resp.Header.Get("Content-Type"), ";")[0]
	if contentType != "" && contentType != "text/html" {
		metadata = append(metadata, contentType)
	}

	if contentType == "text/html" && resp.ContentLength == -1 {
		tree, _ := h5.New(resp.Body)
		titles := getTitles(tree)
		if (len(titles) > 0) {
			metadata = append(metadata, fmt.Sprintf("\x02%s\x02", titles[0].FirstChild.Data))
		}
	} else {
		human := humanSize(resp.ContentLength)
		metadata = append(metadata, human)
	}

	metadata = append(metadata, resp.Request.URL.Host)

	return metadata
}


func getTitles(tree *h5.Tree) (titles []*html.Node) {
	t := transform.New(tree)

	t.Apply(func(node *html.Node) {
		titles = append(titles, node)
	}, "title")

	return titles
}


func humanSize(size int64) string {
	suffixes := []string{"bytes", "KB", "MB", "GB", "TB"}
	for _, suffix := range(suffixes) {
		if size < 1024 {
			return fmt.Sprintf("%d %s", size, suffix)
		}
		size = size/1024
	}
	return "huge"
}
