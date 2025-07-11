package groupdata

import (
	"errors"
	"fmt"
	"pimtrace"
	"strings"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Row struct {
	Headers  map[string]int
	Row      []pimtrace.Value
	Contents pimtrace.Data
}

var _ pimtrace.Entry = (*Row)(nil)
var _ pimtrace.HasStringArray = (*Row)(nil)

type Header interface {
	Get(key string) string
}

func (s *Row) Self() *Row {
	return s
}

func (s *Row) HeadersStringArray() (result []string) {
	result = make([]string, len(s.Headers))
	for h, i := range s.Headers {
		result[i] = h
	}
	return
}

func (s *Row) StringArray(header []string) (result []string) {
	for _, v := range s.Row {
		result = append(result, v.String())
	}
	return
}

func (s *Row) Get(key string) (pimtrace.Value, error) {
	ks := strings.SplitN(key, ".", 2)
	switch ks[0] {
	case "sz", "sized":
		return pimtrace.SimpleIntegerValue(s.Contents.Len()), nil
	case "h", "header", "c", "column":
		ks = ks[1:]
		fallthrough
	default:
		n, ok := s.Headers[ks[0]]
		if ok && len(ks) > 0 {
			return s.Row[n], nil
		}
		var r []pimtrace.Value
		for i := 0; i < s.Contents.Len(); i++ {
			sr := s.Contents.Entry(i)
			if sr == nil {
				continue
			}
			g, err := sr.Get(key)
			if err != nil {
				return nil, fmt.Errorf("group row %w: %s", ErrKeyNotFound, key)
			}
			r = append(r, g)
		}
		return pimtrace.SimpleArrayValue(r), nil
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
