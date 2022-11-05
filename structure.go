package pimtrace

import (
	"fmt"
	mail2 "net/mail"
	"strconv"
	"strings"
	"time"
)

type SimpleStringValue string

func (s SimpleStringValue) Truthy() bool {
	return len(s) > 0
}

func (s SimpleStringValue) Elements() int {
	return 1
}

func (s SimpleStringValue) Length() int {
	return len(s)
}

func (s SimpleStringValue) Array() []Value {
	return []Value{s}
}

func (s SimpleStringValue) StringArray() []string {
	return []string{string(s)}
}

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

type SimpleIntegerValue int

func (s SimpleIntegerValue) Truthy() bool {
	return s != 0
}

func (s SimpleIntegerValue) Elements() int {
	return 1
}

func (s SimpleIntegerValue) Length() int {
	return 1
}

func (s SimpleIntegerValue) Array() []Value {
	return []Value{s}
}

func (s SimpleIntegerValue) StringArray() []string {
	return []string{fmt.Sprint(s)}
}

func (s SimpleIntegerValue) Less(jv Value) bool {
	switch jv := jv.(type) {
	case SimpleIntegerValue:
		return s < jv
	default:
		return s == jv
	}
}

func (s SimpleIntegerValue) Equal(jv Value) bool {
	return s == jv
}

func (s SimpleIntegerValue) Time() *time.Time {
	ut := time.Unix(int64(s), 0)
	return &ut
}

func (s SimpleIntegerValue) Integer() *int {
	si := int(s)
	return &si
}

func (s SimpleIntegerValue) Type() Type {
	return Integer
}

func (s SimpleIntegerValue) String() string {
	return fmt.Sprint(int(s))
}

var _ Value = SimpleIntegerValue(0)

type SimpleArrayValue []Value

func (s SimpleArrayValue) Truthy() bool {
	return len(s) > 0
}

func (s SimpleArrayValue) Elements() int {
	return len(s)
}

func (s SimpleArrayValue) Length() int {
	return len(s)
}

func (s SimpleArrayValue) Array() []Value {
	return s
}

func (s SimpleArrayValue) Less(jv Value) bool {
	switch jv := jv.(type) {
	case SimpleArrayValue:
		for i := 0; i < len(s) || i < len(jv); i++ {
			if len(s) >= i {
				return false
			}
			if len(jv) >= i {
				return true
			}
			if s[i].Less(jv[i]) {
				return true
			}
		}
		return true
	default:
		return false
	}
}

func (s SimpleArrayValue) Equal(jv Value) bool {
	switch jv := jv.(type) {
	case SimpleArrayValue:
		if len(s) != len(jv) {
			return false
		}
		for i := 0; i < len(s); i++ {
			if s[i] != jv[i] {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (s SimpleArrayValue) Time() *time.Time {
	return nil
}

func (s SimpleArrayValue) Integer() *int {
	return nil
}

func (s SimpleArrayValue) Type() Type {
	return Array
}

func (s SimpleArrayValue) String() string {
	return fmt.Sprintf("%#v", s)
}

var _ Value = SimpleArrayValue(nil)

type Type int

const (
	String Type = iota
	Integer
	Array
)

func (t Type) String() string {
	switch t {
	case String:
		return "String"
	case Integer:
		return "Integer"
	case Array:
		return "Array"
	}
	return "unknown"
}

type Value interface {
	Type() Type
	String() string
	Array() []Value
	Time() *time.Time
	Integer() *int
	Equal(jv Value) bool
	Less(jv Value) bool
	Elements() int
	Length() int
	Truthy() bool
}

type Entry interface {
	Get(string) Value
}

type Data interface {
	Len() int
	Entry(n int) Entry
	Truncate(n int) Data
	SetEntry(n int, entry Entry) Data
	NewSelf() Data
}
