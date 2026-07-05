package basic

import (
	"embed"
	"io"
	"text/template"
)

//go:embed help.tmpl
var helpFS embed.FS

var helpTemplate = template.Must(template.ParseFS(helpFS, "help.tmpl"))

func PrintHelp(w io.Writer, t string) error {
	return helpTemplate.Execute(w, struct {
		Type string
	}{
		Type: t,
	})
}
