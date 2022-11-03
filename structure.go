package pimtrace

import (
	mail2 "net/mail"
	"strconv"
	"strings"
	"time"
)

type SimpleStringValue string

func (s SimpleStringValue) Less(jv Value) bool {
	switch jv := jv.(type) {
	case SimpleStringValue:
		return strings.Compare(string(s), string(jv)) < 0
	default:
		return s == jv
	}
}

func (s SimpleStringValue) Equal(jv Value) bool {
	switch jv := jv.(type) {
	case SimpleStringValue:
		return string(s) == string(jv)
	default:
		return s == jv
	}
}

func (s SimpleStringValue) Time() *time.Time {
	t, err := mail2.ParseDate(string(s))
	if err != nil || t.UnixNano() == 0 {
		return nil
	}
	return &t
}

func (s SimpleStringValue) Integer() *int {
	i, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		return nil
	}
	ii := int(i)
	return &ii
}

func (s SimpleStringValue) Type() Type {
	return String
}

func (s SimpleStringValue) String() string {
	return string(s)
}

var _ Value = SimpleStringValue("")

type Type int

const (
	String Type = iota
)

type Value interface {
	Type() Type
	String() string
	Time() *time.Time
	Integer() *int
	Equal(jv Value) bool
	Less(jv Value) bool
}

type Entry interface {
	Get(string) Value
}

type Data interface {
	Len() int
	Entry(n int) Entry
	Truncate(n int) Data
	SetEntry(n int, entry Entry)
}
