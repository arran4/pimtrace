module pimtrace

go 1.18

require (
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/arran4/golang-ical v0.0.0-20220517104411-fd89fefb0182
	github.com/emersion/go-mbox v1.0.3
	github.com/emersion/go-message v0.16.0
	github.com/google/go-cmp v0.5.9
	github.com/olekukonko/tablewriter v0.0.5
)

require (
	github.com/emersion/go-textwrapper v0.0.0-20200911093747-65d896831594 // indirect
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/rivo/uniseg v0.1.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/emersion/go-message v0.16.0 => github.com/arran4/go-message v0.0.0-20221009061333-88d073923c5e
