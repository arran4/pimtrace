package pimtrace

import (
	mail2 "net/mail"
	"strconv"
	"time"
)

type SimpleStringValue string

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
}

type Header interface {
	Get(string) string
}

type Entry[T any] interface {
	Get(string) Value
	Header() Header
	Origin() T
}

type Data[T any] interface {
	Len() int
	Entry(n int) Entry[T]
	Truncate(n int) Data[T]
	SetEntry(n int, entry Entry[T])
	Origin() []T
}
