package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func MarkdownHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")

	filePath := req.URL.Path
	filePath = strings.TrimLeft(filePath, "/render/")
	filePath = strings.TrimRight(filePath, ".html")
	filePath = c.dir + "/" + filePath + ".md"

	file, err := os.Open(filePath)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(filePath + " not found"))
		return
	}
	defer file.Close()

	md, err := ioutil.ReadAll(file)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	data := htmlData{
		Body:  string(markdown.ToHTML(md, parser, renderer)),
		Title: filePath,
	}

	tmpl := template.Must(template.New("page").Parse(htmlScaffold))
	var out bytes.Buffer
	err = tmpl.Execute(&out, data)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("template cound not be rendered"))
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(out.Bytes())
}

type htmlData struct {
	Body  string
	Title string
}

var htmlScaffold = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>{{ .Title }}</title>
</head>
<body>
	{{ .Body }}
</body>
</html>`
