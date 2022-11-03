package csvdata

import (
	"pimtrace"
	"strings"
)

type CSVRow struct {
	Headers map[string]int
	Row     []string
}

type Header interface {
	Get(key string) string
}

func (s *CSVRow) Self() *CSVRow {
	return s
}

func (s *CSVRow) Get(key string) pimtrace.Value {
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

type CSVDataType []*CSVRow

func (p CSVDataType) Truncate(n int) pimtrace.Data[*CSVRow] {
	p = (([]*CSVRow)(p))[:n]
	return p
}

func (p CSVDataType) SetEntry(n int, entry pimtrace.Entry[*CSVRow]) {
	(([]*CSVRow)(p))[n] = entry.Self()
}

func (p CSVDataType) Len() int {
	return len([]*CSVRow(p))
}

func (p CSVDataType) Entry(n int) pimtrace.Entry[*CSVRow] {
	if n >= len([]*CSVRow(p)) || n < 0 {
		return nil
	}
	return ([]*CSVRow(p))[n]
}

func (p CSVDataType) Self() []*CSVRow {
	return []*CSVRow(p)
}

var _ pimtrace.Data[*CSVRow] = CSVDataType(nil)
