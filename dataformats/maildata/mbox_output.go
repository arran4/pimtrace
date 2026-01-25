package maildata

import (
	"fmt"
	"pimtrace"

	"github.com/arran4/go-evaluator"
)

type MBoxOutput struct{}

func (m *MBoxOutput) Execute(d pimtrace.Data, ctx *evaluator.Context) (pimtrace.Data, error) {
	if _, ok := d.(pimtrace.MBoxOutputCapable); !ok {
		return nil, fmt.Errorf("data does not support mbox output")
	}
	return d, nil
}
