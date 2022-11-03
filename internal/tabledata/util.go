package tabledata

import (
	"encoding/csv"
	"fmt"
	"io"
)

func ReadCSV(r io.Reader) (Data, error) {
	header := map[string]int{}
	result := []*Row{}
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
		result = append(result, &Row{
			Headers: header,
			Row:     r,
		})
	}
	return Data(result), nil
}
