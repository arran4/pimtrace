package funcs

import (
	"errors"
	"github.com/araddon/dateparse"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseStrict(t *testing.T) {
	for _, test := range []struct {
		Name   string
		Input  string
		Output string
		Err    error
	}{
		{
			Name:   "Should work",
			Input:  "Wed, 8 Feb 2023 19:00:46 +1100 (AEDT)",
			Output: "2023-02-08 08:00:46 +0000 UTC",
			Err:    nil,
		},
		{
			Name:   "Different form",
			Input:  "Wed,   8 Feb 2023 19:00:46 +1100",
			Output: "2023-02-08 08:00:46 +0000 UTC",
			Err:    nil,
		},
		{
			Name:   "Some sort of date error I got in a google takeout",
			Input:  "Wed,  8 Feb 2023 19:00:46 +1100 (AEDT)",
			Output: "2023-02-08 08:00:46 +0000 UTC",
			Err:    nil,
		},
		{
			Name:   "Month out of range",
			Input:  "FRI, 16 AUG 2013  9:39:51 +1000",
			Output: "2013-08-15 23:39:51 +0000 UTC",
			Err:    nil,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			d, err := dateparse.ParseStrict(test.Input)
			if (err != nil) != (test.Err != nil) || (err != nil && !errors.Is(err, test.Err)) {
				if test.Err == nil {
					t.Errorf("Got error when wanted none: %s", err)
				} else if err == nil {
					t.Errorf("Didn't get an error when we were expecting one: %s", test.Err)
				} else {
					t.Errorf("Got %s expected: %s", err, test.Err)
				}
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(test.Output, d.UTC().String()); diff != "" {
				t.Errorf("Outputs differ: %s", diff)
			}
		})
	}
}
