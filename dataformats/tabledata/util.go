package tabledata

import (
	"encoding/csv"
	"fmt"
	"io"
	"pimtrace"
)

func ReadCSVFile(fType string, fName string) (Data, error) {
	return pimtrace.ReadFileWrapper(fType, fName, ReadCSV)
}
func ReadCSV(r io.Reader, fName string) (Data, error) {
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
		rv := make([]pimtrace.Value, len(r))
		for i, e := range r {
			rv[i] = pimtrace.SimpleStringValue(e)
		}
		result = append(result, &Row{
			Headers: header,
			Row:     rv,
		})
	}
	return Data(result), nil
}
