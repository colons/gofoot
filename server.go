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
	Doc string;
}

type DocContext struct {
	Docs []Doc;
}

func GlobalDocsHandler(w http.ResponseWriter, req *http.Request) {
	context := DocContext{}
	funcMap := template.FuncMap{"Config": GlobalConfig.Global}
	p := markdown.NewParser(&markdown.Extensions{Smart: true})

	for _, command := range(Commands) {
		doc := Doc{}

		if IsArgCommand(command) {
			argCommand := argCommandFor(command)
			doc.Invocation = GlobalConfig.Global("comchar") + strings.Join(argCommand.Args, " ")
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

	http.HandleFunc("/", GlobalDocsHandler)
	http.HandleFunc("/docs.css", ServeCSS)

	err := http.ListenAndServe(GlobalConfig.Global("listen"), nil)
	if err != nil {
		fmt.Printf("Couldn't serve: %s", err)
	}
}
