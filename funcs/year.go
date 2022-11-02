package funcs

import (
	"pimtrace"
	"pimtrace/argparsers/basic"
)

func init() {
	basic.RegisterFunction("year", Year)
}

func Year(d pimtrace.Entry) (pimtrace.Value, error) {

}
