package tabledata

import (
	"encoding/csv"
	"fmt"
	"io"
	"pimtrace"
)

func RowsToData(r ...[]*Row) (result []Data) {
	result = make([]Data, len(r), len(r))
	for i, e := range r {
		result[i] = e
	}
	return
}

func ReadCSV(r io.Reader, fType string, fName string, ops ...any) ([]*Row, error) {
	header := map[string]int{}
	var result []*Row
	cr := csv.NewReader(r)
	for l := 0; ; l++ {
		r, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read: %w", err)
		}
		if r == nil {
			break
		}
		if l == 0 {
			for i, c := range r {
				header[c] = i
			}
			continue
		}
		rv := make([]pimtrace.Value, len(r))
		for i, e := range r {
			rv[i] = pimtrace.SimpleStringValue(e)
		}
		result = append(result, &Row{
			Headers: header,
			Row:     rv,
		})
	}
	return result, nil
}
