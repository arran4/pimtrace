package funcs

import (
	"pimtrace"
	"pimtrace/argparsers/basic"
)

func init() {
	basic.RegisterFunction("count", Count)
}

func Count(d pimtrace.Entry) (pimtrace.Value, error) {

}
