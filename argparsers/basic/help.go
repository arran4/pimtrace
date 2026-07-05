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
	FilterQuery           string
	FilterDesc            string
	TableIntro            string
	TableQuery            string
	TableDesc             string
	SummaryQuery          string
	SummaryDesc           string
	SortQuery             string
	ComboQuery            string
	ColRefNote            string
	NotImplementedNote    string
	CannotConvertBackNote string
}

func PrintHelp(w io.Writer, t string) error {
	data := HelpData{Type: t}
	switch t {
	case "csv":
		data.FilterQuery = "\tfilter not c.job icontains .Software"
		data.FilterDesc = "Which filters each row's `job` column that does not contain `Software`"
		data.TableIntro = "Then you can filter the columns using:"
		data.TableQuery = "\tinto table c.job f.year[c.startdate] f.month[c.startdate]"
		data.TableDesc = ""
		data.SummaryQuery = "\tinto summary c.job calculate f.sum[c.pay] f.count"
		data.SummaryDesc = "Which groups the row on job, then creates a sum of pay and a count of elements grouped"
		data.SortQuery = "\tsort f.year[c.startdate] f.month[c.startdate]"
		data.ComboQuery = "\tfilter c.startdate icontains .Software into summary c.job f.year[c.startdate] calculate f.count filter c.year-startdate eq 2022 sort c.job"
		data.ColRefNote = "- Columns are referred to by `c.`"
		data.NotImplementedNote = ""
		data.CannotConvertBackNote = ""
	case "ical":
		data.FilterQuery = "\tfilter not p.SUMMARY icontains .Report"
		data.FilterDesc = "Which filters all events without SUMMARY containing 'Report'"
		data.TableIntro = "Then you can convert the events to tabular form using:"
		data.TableQuery = "\tinto table p.LOCATION f.year[p.DUE] f.month[p.DUE]"
		data.TableDesc = "Which creates table of locations, and the due year and due month."
		data.SummaryQuery = "\tinto summary p.LOCATION f.year[p.DUE] f.month[p.DUE] calculate f.count"
		data.SummaryDesc = "Which groups the calendar invites based on location, due date year, due date month, then adds a count"
		data.SortQuery = "\tsort f.year[p.DUE] f.month[p.DUE]"
		data.ComboQuery = "\tfilter p.SUMMARY icontains .Report into summary p.LOCATION f.year[p.DUE] f.month[p.DUE] calculate f.count filter c.year-DUE eq 2022 sort p.LOCATION"
		data.ColRefNote = "- Properties are referred to with `p.` once in table form it becomes a column referred to by `c.`"
		data.NotImplementedNote = "- I haven't implemented sub components yet"
		data.CannotConvertBackNote = "- Once converted to table form, it can not be converted back to ical / ics"
	case "mail":
		data.FilterQuery = "\tfilter not h.user-agent icontains .Kmail"
		data.FilterDesc = "Which filters all emails sent with a user agent that does not contain `kmail`"
		data.TableIntro = "Then you can convert the emails to tabular form using:"
		data.TableQuery = "\tinto table h.user-agent h.subject f.year[h.date] f.month[h.date]"
		data.TableDesc = "Which creates table of user agent, subject, and the year and month."
		data.SummaryQuery = "\tinto summary h.user-agent h.subject f.year[h.date] f.month[h.date] calculate f.sum[c.size] f.count"
		data.SummaryDesc = "Which groups the mail based on user-agent, subject, year, month, then creates a sum and a count"
		data.SortQuery = "\tsort f.year[h.date] f.month[h.date]"
		data.ComboQuery = "\tfilter h.user-agent icontains .Kmail into summary h.user-agent f.year[h.date] calculate f.count filter c.year-date eq 2022 sort h.user-agent"
		data.ColRefNote = "- Headers are referred to with `h.` once in table form it becomes a column referred to by `c.`"
		data.NotImplementedNote = "- I haven't implemented any body parsing components yet"
		data.CannotConvertBackNote = "- Once converted to table form, it can not be converted back to mailfile or mbox"
	}
	return helpTemplate.Execute(w, data)
}
