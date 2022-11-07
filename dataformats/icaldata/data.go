package icaldata

import (
	"errors"
	"fmt"
	"github.com/arran4/golang-ical"
	"pimtrace"
	"strings"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrHeaderError = errors.New("header error")
)

type ICalWithSource struct {
	Component     ics.Component
	ComponentBase *ics.ComponentBase
	SourceType    string
	SourceFile    string
	Header        map[string]int
}

var _ pimtrace.Entry = (*ICalWithSource)(nil)
var _ pimtrace.HasStringArray = (*ICalWithSource)(nil)

func (s *ICalWithSource) Self() *ICalWithSource {
	return s
}

func (s *ICalWithSource) HeadersStringArray() (result []string) {
	result = make([]string, 0, len(s.Header))
	for h := range s.Header {
		result = append(result, h)
	}
	return
}

func (s *ICalWithSource) StringArray(header []string) (result []string) {
	for _, h := range header {
		i, ok := s.Header[h]
		if !ok {
			continue
		}
		result = append(result, s.ComponentBase.Properties[i].Value)
	}
	return
}

func (s *ICalWithSource) Get(key string) (pimtrace.Value, error) {
	ks := strings.SplitN(key, ".", 2)
	switch ks[0] {
	//case "sz", "sized": TODO
	//	return SimpleNumberValue(s.
	case "p", "property":
		ks = ks[1:]
		fallthrough
	default:
		if len(ks) > 1 {
			i, ok := s.Header[ks[0]]
			if !ok {
				return nil, fmt.Errorf("iCal get %w, %s", ErrHeaderError, key)
			}
			return pimtrace.SimpleStringValue(s.ComponentBase.Properties[i].Value), nil
		}
		return nil, fmt.Errorf("iCal get %w, %s", ErrKeyNotFound, key)
	}
}

type Data []*ICalWithSource

func (icd Data) Truncate(n int) pimtrace.Data {
	icd = (([]*ICalWithSource)(icd))[:n]
	return icd
}

func (icd Data) SetEntry(n int, entry pimtrace.Entry) pimtrace.Data {
	for n > len(icd) {
		icd = append((([]*ICalWithSource)(icd)), nil)
	}
	if n == len(icd) {
		icd = append(icd, entry.(*ICalWithSource))
	} else {
		(([]*ICalWithSource)(icd))[n] = entry.(*ICalWithSource)
	}
	return icd

}

func (icd Data) Len() int {
	return len([]*ICalWithSource(icd))
}

func (icd Data) Entry(n int) pimtrace.Entry {
	if n >= len([]*ICalWithSource(icd)) || n < 0 {
		return nil
	}
	return ([]*ICalWithSource(icd))[n]
}

func (icd Data) Self() []*ICalWithSource {
	return []*ICalWithSource(icd)
}

func (icd Data) NewSelf() pimtrace.Data {
	return Data(make([]*ICalWithSource, 0))
}

var _ pimtrace.Data = Data(nil)
