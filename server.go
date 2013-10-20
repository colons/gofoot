package main

import (
	"fmt"
	"net/http"
	"github.com/knieriem/markdown"
	"text/template"
	"strings"
)

var GlobalConfig config
var NetworkConfigs map[string]config
var MarkdownParser *markdown.Parser
var Template *template.Template

type Doc struct {
	Invocation string
	Doc string;
}

type DocContext struct {
	Docs []Doc;
}

func GlobalDocsHandler(w http.ResponseWriter, req *http.Request) {
	context := DocContext{}
	funcMap := template.FuncMap{"Config": GlobalConfig.Global}

	for _, command := range(Commands) {
		doc := Doc{}

		if IsArgCommand(command) {
			argCommand := argCommandFor(command)
			doc.Invocation = GlobalConfig.Global("comchar") + strings.Join(argCommand.Args, " ")
		}
		doc.Doc = command.GetDocs()

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
	MarkdownParser = markdown.NewParser(&markdown.Extensions{Smart: true})

	http.HandleFunc("/", GlobalDocsHandler)
	http.HandleFunc("/docs.css", ServeCSS)

	err := http.ListenAndServe(GlobalConfig.Global("listen"), nil)
	if err != nil {
		fmt.Printf("Couldn't serve: %s", err)
	}
}
