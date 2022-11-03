package tabledata

import (
	"fmt"
	"os"
	"pimtrace"
	"strings"
)

type Row struct {
	Headers map[string]int
	Row     []string
}

type Header interface {
	Get(key string) string
}

func (s *Row) Self() *Row {
	return s
}

func (s *Row) Get(key string) pimtrace.Value {
	ks := strings.SplitN(key, ".", 2)
	switch ks[0] {
	//case "sz", "sized": TODO
	//	return SimpleNumberValue(s.
	case "h", "header", "c", "column":
		fallthrough
	default:
		n, ok := s.Headers[ks[1]]
		if ok && len(ks) > 1 {
			return pimtrace.SimpleStringValue(s.Row[n])
		}
		return nil
	}
}

type Data []*Row

func (p Data) Truncate(n int) pimtrace.Data[*Row] {
	p = (([]*Row)(p))[:n]
	return p
}

func (p Data) SetEntry(n int, entry pimtrace.Entry[*Row]) {
	(([]*Row)(p))[n] = entry.Self()
}

func (p Data) Len() int {
	return len([]*Row(p))
}

func (p Data) Entry(n int) pimtrace.Entry[*Row] {
	if n >= len([]*Row(p)) || n < 0 {
		return nil
	}
	return ([]*Row(p))[n]
}

func (p Data) Self() []*Row {
	return []*Row(p)
}

func (p Data) Output(mode, outputPath string) error {
	switch mode {
	case "csv":
		switch outputPath {
		case "-":
			return WriteCSVStream(p, os.Stdin, outputPath)
		default:
			return WriteCSVFile(p, outputPath)
		}
	case "count":
		fmt.Println(p.Len())
		return nil
	case "list":
		fmt.Println("`--output-type`s: ")
		fmt.Printf(" =%-20s - %s\n", "list", "This help text")
		fmt.Printf(" =%-20s - %s\n", "count", "Just a count")
		fmt.Printf(" =%-20s - %s\n", "csv", "Data in csv format")
		fmt.Println()
		return nil
	default:
		//fmt.Println("Please specify a -input-type")
		//fmt.Println()
		return nil
	}
}

var _ pimtrace.Data[*Row] = Data(nil)
