package nildata

import (
	"pimtrace"
)

type Row struct {
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
	return
}

func (s *Row) StringArray(header []string) (result []string) {
	return
}

func (s *Row) Get(key string) (pimtrace.Value, error) {
	return &pimtrace.SimpleNilValue{}, nil
}

type Data []*Row

func (d Data) Truncate(n int) pimtrace.Data {
	return d
}

func (d Data) SetEntry(n int, entry pimtrace.Entry) pimtrace.Data {
	return d
}

func (d Data) Len() int {
	return 0
}

func (d Data) Entry(n int) pimtrace.Entry {
	return &Row{}
}

func (d Data) Self() []*Row {
	return []*Row(d)
}

func (d Data) NewSelf() pimtrace.Data {
	return Data(make([]*Row, 0))
}

var _ pimtrace.Data = Data(nil)
