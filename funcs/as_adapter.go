package funcs

import (
	"fmt"
)

type AsAdapter struct{}

func (a *AsAdapter) Call(args ...interface{}) (interface{}, error) {
	// AsAdapter in expression context might just return the value?
	// Or a tuple?
	// If used in 'select f.as(val, name)', it might return val?
	// The naming happens in AST usually.
	// If executed, we return the first argument.
	if len(args) == 0 {
		return nil, fmt.Errorf("expected at least 1 argument")
	}
	return args[0], nil
}
