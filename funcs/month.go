package funcs

import (
	"pimtrace"
	"pimtrace/argparsers/basic"
)

func init() {
	basic.RegisterFunction("month", Month)
}

func Month(d pimtrace.Entry) (pimtrace.Value, error) {

}
