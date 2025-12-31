package maildata

import (
	"fmt"
	"pimtrace"
)

type MBoxOutput struct{}

func (m *MBoxOutput) Execute(d pimtrace.Data) (pimtrace.Data, error) {
	if _, ok := d.(pimtrace.MBoxOutputCapable); !ok {
		return nil, fmt.Errorf("data does not support mbox output")
	}
	return d, nil
}
