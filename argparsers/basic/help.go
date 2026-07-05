package basic

import (
	"embed"
	"io"
	"text/template"
)

//go:embed help.tmpl
var helpFS embed.FS

var helpTemplate = template.Must(template.ParseFS(helpFS, "help.tmpl"))

type HelpData struct {
	Type                  string
	Entity                string
	EntityPlural          string
	FilterExample         string
	FilterExplanation     string
	TableExample          string
	SummaryExample        string
	SortExample           string
	ComboExample          string
	ColRef                string
	NotImplementedNote    string
	CannotConvertBackNote string
}

func PrintHelp(w io.Writer, t string) error {
	data := HelpData{Type: t}
	switch t {
	case "csv":
		data.Entity = "row"
		data.EntityPlural = "rows"
		data.FilterExample = "filter not c.job icontains .Software"
		data.FilterExplanation = "Filters all rows where the `job` column does not contain `Software`."
		data.TableExample = "into table c.job f.year[c.startdate] f.month[c.startdate]"
		data.SummaryExample = "into summary c.job calculate f.sum[c.pay] f.count"
		data.SortExample = "sort f.year[c.startdate] f.month[c.startdate]"
		data.ComboExample = "filter c.startdate icontains .Software into summary c.job f.year[c.startdate] calculate f.count filter c.year-startdate eq 2022 sort c.job"
		data.ColRef = "Columns are referred to by `c.`."
	case "ical":
		data.Entity = "event"
		data.EntityPlural = "events"
		data.FilterExample = "filter not p.SUMMARY icontains .Report"
		data.FilterExplanation = "Filters all events without a SUMMARY containing 'Report'."
		data.TableExample = "into table p.LOCATION f.year[p.DUE] f.month[p.DUE]"
		data.SummaryExample = "into summary p.LOCATION f.year[p.DUE] f.month[p.DUE] calculate f.count"
		data.SortExample = "sort f.year[p.DUE] f.month[p.DUE]"
		data.ComboExample = "filter p.SUMMARY icontains .Report into summary p.LOCATION f.year[p.DUE] f.month[p.DUE] calculate f.count filter c.year-DUE eq 2022 sort p.LOCATION"
		data.ColRef = "Properties are referred to with `p.`. Once in table form, they become columns referred to by `c.`."
		data.NotImplementedNote = "- Subcomponents are not fully implemented yet."
		data.CannotConvertBackNote = "- Once converted to table form, the data cannot be converted back to ical / ics."
	case "mail":
		data.Entity = "email"
		data.EntityPlural = "emails"
		data.FilterExample = "filter not h.user-agent icontains .Kmail"
		data.FilterExplanation = "Filters all emails sent with a user agent that does not contain `kmail`."
		data.TableExample = "into table h.user-agent h.subject f.year[h.date] f.month[h.date]"
		data.SummaryExample = "into summary h.user-agent h.subject f.year[h.date] f.month[h.date] calculate f.sum[c.size] f.count"
		data.SortExample = "sort f.year[h.date] f.month[h.date]"
		data.ComboExample = "filter h.user-agent icontains .Kmail into summary h.user-agent f.year[h.date] calculate f.count filter c.year-date eq 2022 sort h.user-agent"
		data.ColRef = "Headers are referred to with `h.`. Once in table form, they become columns referred to by `c.`."
		data.NotImplementedNote = "- Body parsing components are not fully implemented yet."
		data.CannotConvertBackNote = "- Once converted to table form, the data cannot be converted back to mailfile or mbox."
	}
	return helpTemplate.Execute(w, data)
}
