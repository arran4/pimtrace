package dataformats

import (
	"fmt"
	"io"
	"pimtrace"
	"pimtrace/ast"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/arran4/lookup"
)

type FieldDescription struct {
	Name          string
	Expression    string
	Component     string
	ObservedTypes map[string]bool
	SampleCount   int
	MatchCount    int
}

func Describe(data pimtrace.Data, maxSample int) []FieldDescription {
	if data == nil || data.Len() == 0 {
		return nil
	}

	n := data.Len()
	if maxSample > 0 && n > maxSample {
		n = maxSample
	}

	fieldMap := make(map[string]*FieldDescription)
	var orderedKeys []string

	for i := 0; i < n; i++ {
		entry := data.Entry(i)
		if entry == nil {
			continue
		}

		var fields []pimtrace.FieldDescriptor

		if sd, ok := entry.(pimtrace.SchemaDescriber); ok {
			fields = sd.SchemaFields()
		} else if hsa, ok := entry.(pimtrace.HasStringArray); ok {
			// Fallback
			headers := hsa.HeadersStringArray()
			for _, h := range headers {
				fields = append(fields, pimtrace.FieldDescriptor{
					Name:            h,
					QueryExpression: "c." + h,
					Component:       "unknown",
				})
			}
		}

		for _, field := range fields {
			if !pimtrace.IsQueryable(field.Name) {
				continue
			}
			expr := field.QueryExpression
			comp := field.Component

			key := expr + "|" + comp
			if _, exists := fieldMap[key]; !exists {
				fieldMap[key] = &FieldDescription{
					Name:          field.Name,
					Expression:    expr,
					Component:     comp,
					ObservedTypes: make(map[string]bool),
					SampleCount:   n,
				}
				orderedKeys = append(orderedKeys, key)
			}

			desc := fieldMap[key]

			// Evaluate expression via standard entry lookup to ensure accuracy
			ep := ast.NewEntryPathor(entry)
			valPathor := ep.Find(expr)

			var rawVal interface{}
			var resolveErr error

			if _, ok := valPathor.(*lookup.Invalidor); ok {
				resolveErr = fmt.Errorf("invalid path")
			} else if valPathor != nil {
				rawVal = valPathor.Raw()
				if errRaw, isErr := rawVal.(error); isErr {
					resolveErr = errRaw
				}
			} else {
				resolveErr = fmt.Errorf("nil pathor")
			}

			if resolveErr == nil {
				desc.MatchCount++
				typeStr := "nil"
				if rawVal != nil {
					if pimVal, ok := rawVal.(pimtrace.Value); ok {
						typeStr = strings.ToLower(fmt.Sprintf("%v", pimVal.Type()))
					} else {
						t := reflect.TypeOf(rawVal)
						if t != nil {
							name := t.Name()
							if name == "" {
								name = t.String()
							}
							typeStr = strings.ToLower(strings.TrimPrefix(name, "pimtrace."))
						} else {
							typeStr = "unknown"
						}
					}
				}
				desc.ObservedTypes[typeStr] = true
			}
		}
	}

	var results []FieldDescription
	for _, k := range orderedKeys {
		if fieldMap[k].MatchCount > 0 {
			results = append(results, *fieldMap[k])
		}
	}

	// Deterministic sort by Name, then Expression, then Component
	sort.Slice(results, func(i, j int) bool {
		if results[i].Name == results[j].Name {
			if results[i].Expression == results[j].Expression {
				return results[i].Component < results[j].Component
			}
			return results[i].Expression < results[j].Expression
		}
		return results[i].Name < results[j].Name
	})

	return results
}

func DescribeAndPrint(data pimtrace.Data, w io.Writer) error {
	sampleBound := 1000
	descs := Describe(data, sampleBound)
	if len(descs) == 0 {
		_, err := fmt.Fprintln(w, "No queryable fields discovered (or data is empty).")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	_, _ = fmt.Fprintln(w, "Note: Values are sampled up to a maximum of 1000 records to infer types.")
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(tw, "Component\tName\tExpression\tObserved Types\tMatches (Sampled)")
	_, _ = fmt.Fprintln(tw, "---------\t----\t----------\t--------------\t-----------------")

	for _, d := range descs {
		var types []string
		for t := range d.ObservedTypes {
			types = append(types, t)
		}
		sort.Strings(types)
		typeStr := strings.Join(types, ", ")
		if typeStr == "" {
			typeStr = "absent"
		}

		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d/%d\n", d.Component, d.Name, d.Expression, typeStr, d.MatchCount, d.SampleCount)
	}

	return tw.Flush()
}
