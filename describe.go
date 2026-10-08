package pimtrace

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"

	)

type FieldDescriptor struct {
	Name            string
	QueryExpression string
	Component       string
}

type SchemaDescriber interface {
	SchemaFields() []FieldDescriptor
}

// IsQueryable checks if a field name can be safely used in the current basic parser
// without introducing syntax errors or unintended token splitting.
func IsQueryable(s string) bool {
	if s == "" {
		return false
	}
	if strings.ContainsAny(s, " \t\n\r()\"'=<>!&|[]") {
		return false
	}
	return true
}

type FieldDescription struct {
	Name         string
	Expression   string
	Component    string
	ObservedType string
	Count        int
}

func Describe(data Data, maxSample int) []FieldDescription {
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

		var fields []FieldDescriptor

		if sd, ok := entry.(SchemaDescriber); ok {
			fields = sd.SchemaFields()
		} else if hsa, ok := entry.(HasStringArray); ok {
			// Fallback for non-compliant structs
			headers := hsa.HeadersStringArray()
			for _, h := range headers {
				fields = append(fields, FieldDescriptor{
					Name:            h,
					QueryExpression: "c." + h,
					Component:       "unknown",
				})
			}
		}

		for _, field := range fields {
			if !IsQueryable(field.Name) {
				continue
			}
			expr := field.QueryExpression

			key := expr
			if _, exists := fieldMap[key]; !exists {
				fieldMap[key] = &FieldDescription{
					Name:       field.Name,
					Expression: expr,
					Component:  field.Component,
				}
				orderedKeys = append(orderedKeys, key)
			}

			desc := fieldMap[key]
			desc.Count++

			if desc.ObservedType == "" || desc.ObservedType == "nil" {
				rawVal, _ := entry.Get(expr)

				if rawVal != nil {
					switch rawVal.(type) {
					case *SimpleStringValue:
						desc.ObservedType = "string"
					case *SimpleIntegerValue:
						desc.ObservedType = "integer"
					case *SimpleFloatValue:
						desc.ObservedType = "float"
					case *SimpleNilValue:
						desc.ObservedType = "nil"
					default:
						t := reflect.TypeOf(rawVal)
						if t != nil {
							name := t.Name()
							if name == "" {
								name = t.String()
							}
							desc.ObservedType = strings.TrimPrefix(name, "pimtrace.")
						} else {
							desc.ObservedType = "unknown"
						}
					}
				} else {
					desc.ObservedType = "nil"
				}
			}
		}
	}

	var results []FieldDescription
	for _, k := range orderedKeys {
		results = append(results, *fieldMap[k])
	}
	return results
}

func DescribeAndPrint(data Data, w io.Writer) error {
	descs := Describe(data, 1000)
	if len(descs) == 0 {
		_, err := fmt.Fprintln(w, "No queryable fields discovered (or data is empty).")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	_, _ = fmt.Fprintln(tw, "Name\tExpression\tType\tCount\t")
	_, _ = fmt.Fprintln(tw, "----\t----------\t----\t-----\t")

	for _, d := range descs {
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t\n", d.Name, d.Expression, d.ObservedType, d.Count)
	}

	return tw.Flush()
}
