package ast

import (
	"pimtrace"

	"github.com/arran4/lookup"
)

type EntryPathor struct {
	lookup.Pathor
	Entry pimtrace.Entry
}

var _ lookup.Finder = (*EntryPathor)(nil)
var _ lookup.Pathor = (*EntryPathor)(nil)

func NewEntryPathor(e pimtrace.Entry) *EntryPathor {
	return &EntryPathor{
		Pathor: lookup.Reflect(e),
		Entry:  e,
	}
}

func (e *EntryPathor) Find(path string, opts ...lookup.Runner) lookup.Pathor {
	// Try standard lookup first?
	// Or try Entry.Get first?
	// The goal is to use Entry.Get for "c.date" etc.
	// But lookup.Reflector.Find splits path?
	// If I pass "c.date", Reflector might try Find("c").Find("date").
	// But Entry.Get expects "c.date".

	// So I should try Entry.Get(path).
	// If path is empty, return self.
	if path == "" {
		return e
	}

	val, err := e.Entry.Get(path)
	if err == nil {
		// Found it!
		// Return a Pathor for the result.
		// If result is an Entry, wrap it again?
		if subEntry, ok := val.(pimtrace.Entry); ok {
			return NewEntryPathor(subEntry)
		}
		return lookup.Reflect(val)
	}

	// If Entry.Get fails, maybe fallback to standard reflection?
	// But Entry acts as the source of truth.
	// Return Invalidor.
	return lookup.NewInvalidor(path, err)
}
