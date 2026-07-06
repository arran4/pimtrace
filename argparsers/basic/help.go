package basic

import (
	"embed"
	"io"
	"text/template"
)

//go:embed help_*.tmpl
var helpFS embed.FS

var helpTemplate = template.Must(template.ParseFS(helpFS, "help_*.tmpl"))

func PrintHelp(w io.Writer, t string) error {
	return helpTemplate.ExecuteTemplate(w, "help_"+t+".tmpl", nil)
}
