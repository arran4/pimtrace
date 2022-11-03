package groupdata

import (
	"pimtrace"
	"strings"
)

type Row struct {
	Headers  map[string]int
	Row      []pimtrace.Value
	Contents pimtrace.Data
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
			return s.Row[n]
		}
		var r []pimtrace.Value
		for i := 0; i < s.Contents.Len(); i++ {
			sr := s.Contents.Entry(i)
			if sr == nil {
				continue
			}
			r = append(r, sr.Get(key))
		}
		return pimtrace.SimpleArrayValue(r)
	}
}

type Data []*Row

func (d Data) Truncate(n int) pimtrace.Data {
	d = (([]*Row)(d))[:n]
	return d
}

func (d Data) SetEntry(n int, entry pimtrace.Entry) pimtrace.Data {
	for n > len(d) {
		d = append((([]*Row)(d)), nil)
	}
	if n == len(d) {
		d = append(d, entry.(*Row))
	} else {
		(([]*Row)(d))[n] = entry.(*Row)
	}
	return d
}

func (d Data) Len() int {
	return len([]*Row(d))
}

func (d Data) Entry(n int) pimtrace.Entry {
	if n >= len([]*Row(d)) || n < 0 {
		return nil
	}
	return ([]*Row(d))[n]
}

func (d Data) Self() []*Row {
	return []*Row(d)
}

func (d Data) NewSelf() pimtrace.Data {
	return Data(make([]*Row, 0))
}

var _ pimtrace.Data = Data(nil)
