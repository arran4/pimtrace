package funcs

import (
	"pimtrace"
	"pimtrace/argparsers/basic"
)

func init() {
	basic.RegisterFunction("sum", Sum)
}

func Sum(d pimtrace.Entry) (pimtrace.Value, error) {

}
