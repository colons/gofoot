package main

import (
	"fmt"
	"net/http"
	"github.com/knieriem/markdown"
	"text/template"
	"strings"
	"bytes"
)

var NetworkConfigs map[string]config
var Template *template.Template

type Doc struct {
	Invocation string
	Doc string
}

type DocContext struct {
	Docs []Doc
}

func HelpRedirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/help/", 302)
}

func DocsHandler(w http.ResponseWriter, req *http.Request) {
	slashlessPath := strings.Trim(req.URL.Path, "/")
	pathElements := strings.Split(slashlessPath, "/")[1:]
	var config func(string) string

	switch elementCount := len(pathElements); elementCount {
	case 0:
		config = GlobalConfig.Global
	case 1:
		config = GetConfig(pathElements[0]).Network
	case 2:
		config = func(key string) string {
			thisConfig := GetConfig(pathElements[0])
			return thisConfig.Source(pathElements[1], key)
		}
	default:
		fmt.Println(elementCount)
		config = GlobalConfig.Global
	}

	context := DocContext{}
	funcMap := template.FuncMap{"config": config}
	p := markdown.NewParser(&markdown.Extensions{Smart: true})

	for _, command := range(Commands) {
		if !commandNotIn(command, config("blacklist")) {
			continue
		}

		doc := Doc{}

		if IsArgCommand(command) {
			argCommand := argCommandFor(command)
			doc.Invocation = config("comchar") + strings.Join(argCommand.Args, " ")
		}

		mdDocBuffer := bytes.NewBufferString(command.GetDocs())
		htmlDocBuffer := bytes.NewBuffer([]byte{})
		p.Markdown(mdDocBuffer, markdown.ToHTML(htmlDocBuffer))
		doc.Doc = htmlDocBuffer.String()

		context.Docs = append(context.Docs, doc)
	}

	ourTemplate := template.New("docs.html").Funcs(funcMap)
	ourTemplate = template.Must(ourTemplate.ParseFiles("docs/docs.html"))
	err := ourTemplate.Execute(w, context)

	if err != nil {
		fmt.Printf("Error rendering: %s\n", err)
	}
}

func ServeCSS(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "docs/docs.css")
}

func RunServer() {
	GlobalConfig = GetConfig("")

	http.HandleFunc("/", HelpRedirect)
	http.HandleFunc("/help/", DocsHandler)
	http.HandleFunc("/docs.css", ServeCSS)

	err := http.ListenAndServe(GlobalConfig.Global("listen"), nil)
	if err != nil {
		fmt.Printf("Couldn't serve: %s", err)
	}
}
