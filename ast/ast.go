package ast

import (
	"fmt"
	"pimtrace"
	"pimtrace/dataformats/groupdata"
	"pimtrace/dataformats/tabledata"
	"pimtrace/funcs"
	"sort"
	"strings"
	"unicode"

	"github.com/arran4/go-evaluator"
	"github.com/arran4/lookup"
)

var (
	ErrUnknownFunction = fmt.Errorf("unknown function")
)

type Operation interface {
	Execute(d pimtrace.Data) (pimtrace.Data, error)
}

type CompoundStatement struct {
	Statements []Operation
}

func (o *CompoundStatement) Execute(d pimtrace.Data) (pimtrace.Data, error) {
	for _, op := range o.Statements {
		var err error
		d, err = op.Execute(d)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (o *CompoundStatement) Simplify() Operation {
	if len(o.Statements) == 0 {
		return nil
	}
	if len(o.Statements) == 1 {
		return o.Statements[0]
	}
	var result []Operation
	for i, statement := range o.Statements {
		switch statement := statement.(type) {
		case *CompoundStatement:
			if result == nil {
				result = append([]Operation{}, o.Statements[:i]...)
			}
			result = append(result, statement.Statements...)
		default:
			if result == nil {
				continue
			}
			result = append(result, statement)
		}
	}
	if result != nil {
		o.Statements = result
	}
	return o
}

var _ Operation = (*CompoundStatement)(nil)

type ValueExpression interface {
	funcs.ValueExpression
	ColumnName() string
	evaluator.Term
}

type ConstantExpression string

func (ve ConstantExpression) ColumnName() string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return '-'
	}, string(ve))
}

func (ve ConstantExpression) Execute(d pimtrace.Entry) (pimtrace.Value, error) {
	return pimtrace.SimpleStringValue(ve), nil
}

func (ve ConstantExpression) Evaluate(d interface{}) (interface{}, error) {
	return pimtrace.SimpleStringValue(ve), nil
}

type EntryExpression string

func (ve EntryExpression) ColumnName() string {
	ss := strings.SplitN(string(ve), ".", 2)
	s := ""
	if len(ss) > 1 {
		s = ss[1]
	}
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return '-'
	}, s)
}

func (ve EntryExpression) Execute(d pimtrace.Entry) (pimtrace.Value, error) {
	return d.Get(string(ve))
}

func (ve EntryExpression) Evaluate(d interface{}) (interface{}, error) {
	if w, ok := d.(evaluatorEntryWrapper); ok {
		d = w.Entry
	}
	eEntry, ok := d.(pimtrace.Entry)
	if !ok {
		return nil, fmt.Errorf("invalid entry type")
	}

	// Create adapter
	ep := NewEntryPathor(eEntry)

	// Use lookup to traverse
	// We need to pass the path string(ve) which might be "c.date"
	// lookup.Reflect(ep) returns a Reflector wrapping EntryPathor.
	// Reflector.Find(path) will call EntryPathor.Find(path) because we implemented Finder interface check in Reflector.

	res := lookup.Reflect(ep).Find(string(ve))
	return res.Raw(), nil
}

type OpFunc func(pimtrace.Value, pimtrace.Value) (bool, error)

func EqualOp(rhsv pimtrace.Value, lhsv pimtrace.Value) (bool, error) {
	return rhsv.String() == lhsv.String(), nil
}

var _ OpFunc = EqualOp

func ContainsOp(rhsv pimtrace.Value, lhsv pimtrace.Value) (bool, error) {
	return strings.Contains(rhsv.String(), lhsv.String()), nil
}

var _ OpFunc = ContainsOp

func IContainsOp(rhsv pimtrace.Value, lhsv pimtrace.Value) (bool, error) {
	return strings.Contains(strings.ToLower(rhsv.String()), strings.ToLower(lhsv.String())), nil
}

var _ OpFunc = IContainsOp

type Op struct {
	Op  string
	LHS ValueExpression
	RHS ValueExpression
}

func (e *Op) Evaluate(d interface{}) bool {
	if e.LHS == nil || e.RHS == nil {
		return false
	}

	if w, ok := d.(evaluatorEntryWrapper); ok {
		d = w.Entry
	}
	eEntry, ok := d.(pimtrace.Entry)
	if !ok {
		return false
	}

	expr := evaluator.ComparisonExpression{
		LHS:       evaluator.Self{},
		RHS:       evaluator.Self{},
		Operation: e.Op,
	}

	expr.LHS = e.LHS
	expr.RHS = e.RHS

	return expr.Evaluate(eEntry)
}

type FilterStatement struct {
	Expression *evaluator.Query
}

func (f FilterStatement) Execute(d pimtrace.Data) (pimtrace.Data, error) {
	return Filter(d, f.Expression)
}

var _ Operation = (*FilterStatement)(nil)

type FunctionExpression struct {
	Function string
	Args     []ValueExpression
	F        funcs.Function[ValueExpression]
}

func (fe *FunctionExpression) ColumnName() string {
	fe.LoadFunction()
	switch f := fe.F.(type) {
	case funcs.ColumnNamer[ValueExpression]:
		v := f.ColumnName(fe.Args)
		if len(v) > 0 {
			return v
		}
	}
	elems := []string{fe.Function}
	for _, arg := range fe.Args {
		elems = append(elems, arg.ColumnName())
	}
	return strings.Join(elems, "-")
}

func (fe *FunctionExpression) Execute(d pimtrace.Entry) (pimtrace.Value, error) {
	fe.LoadFunction()
	if fe.F == nil {
		return nil, fmt.Errorf("%w: %s", ErrUnknownFunction, fe.Function)
	}
	return fe.F.Run(d, fe.Args)
}

func (fe *FunctionExpression) LoadFunction() {
	if fe.F == nil {
		functions := funcs.Functions[ValueExpression]()
		if f, ok := functions[fe.Function]; ok {
			fe.F = f
		}
	}
}

func (fe *FunctionExpression) Evaluate(d interface{}) (interface{}, error) {
	if w, ok := d.(evaluatorEntryWrapper); ok {
		d = w.Entry
	}
	eEntry, ok := d.(pimtrace.Entry)
	if !ok {
		return nil, fmt.Errorf("invalid entry type")
	}
	return fe.Execute(eEntry)
}

var _ ValueExpression = (*FunctionExpression)(nil)

type EvaluatorFunctionExpression struct {
	Function string
	evaluator.FunctionExpression
}

func (fe *EvaluatorFunctionExpression) ColumnName() string {
	elems := []string{fe.Function}
	for _, arg := range fe.Args {
		if c, ok := arg.(ValueExpression); ok {
			elems = append(elems, c.ColumnName())
		}
	}
	return strings.Join(elems, "-")
}

func (fe *EvaluatorFunctionExpression) Execute(d pimtrace.Entry) (pimtrace.Value, error) {
	res, err := fe.Evaluate(d)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(pimtrace.Value); ok {
		return v, nil
	}
	// Conversion logic if result is not pimtrace.Value
	return toPimtraceValue(res)
}

func toPimtraceValue(v interface{}) (pimtrace.Value, error) {
	if v == nil {
		return &pimtrace.SimpleNilValue{}, nil
	}
	switch val := v.(type) {
	case pimtrace.Value:
		return val, nil
	case int:
		return pimtrace.SimpleIntegerValue(val), nil
	case int64:
		return pimtrace.SimpleIntegerValue(int(val)), nil
	case float64:
		return pimtrace.SimpleIntegerValue(int(val)), nil // Lossy? pimtrace seems to use int mostly
	case string:
		return pimtrace.SimpleStringValue(val), nil
	case []interface{}:
		var arr []pimtrace.Value
		for _, item := range val {
			pv, err := toPimtraceValue(item)
			if err != nil {
				return nil, err
			}
			arr = append(arr, pv)
		}
		return pimtrace.SimpleArrayValue(arr), nil
	}
	return nil, fmt.Errorf("cannot convert %T to pimtrace.Value", v)
}

// Evaluate is inherited from evaluator.FunctionExpression but we might need to wrap context?
// No, FunctionExpression.Evaluate does: arg.Evaluate(d).
// As long as arguments are ValueExpressions (which trigger Finder), it works.
// BUT arguments stored in evaluator.FunctionExpression are []Term.
// ValueExpression implements Term. So it works.

var _ ValueExpression = (*EvaluatorFunctionExpression)(nil)

type FunctionDef func(d pimtrace.Entry) (pimtrace.Value, error)

type ColumnExpression struct {
	Name      string
	Operation ValueExpression
}

type TableTransformer struct {
	Columns []*ColumnExpression
}

func (t *TableTransformer) Execute(d pimtrace.Data) (pimtrace.Data, error) {
	headers := map[string]int{}
	for i, c := range t.Columns {
		headers[c.Name] = i
	}
	td := make([]*tabledata.Row, d.Len())
	for i := 0; i < d.Len(); i++ {
		r := make([]pimtrace.Value, len(t.Columns))
		e := d.Entry(i)
		for i, c := range t.Columns {
			v, err := c.Operation.Execute(e)
			if err != nil {
				return nil, err
			}
			r[i] = v
		}
		td[i] = &tabledata.Row{
			Headers: headers,
			Row:     r,
		}
	}
	return tabledata.Data(td), nil
}

var _ Operation = (*TableTransformer)(nil)

type SortTransformer struct {
	Expression []ValueExpression
}

type SortTransformerSorter struct {
	SortTransformer *SortTransformer
	Data            pimtrace.Data
}

func (s *SortTransformerSorter) Len() int {
	return s.Data.Len()
}

func (s *SortTransformerSorter) Less(i, j int) bool {
	for _, e := range s.SortTransformer.Expression {
		io, jo := s.Data.Entry(i), s.Data.Entry(j)
		iv, _ := e.Execute(io)
		jv, _ := e.Execute(jo)
		if iv.Equal(jv) {
			continue
		}
		return iv.Less(jv)
	}
	return i < j
}

func (s *SortTransformerSorter) Swap(i, j int) {
	io, jo := s.Data.Entry(i), s.Data.Entry(j)
	s.Data.SetEntry(j, io)
	s.Data.SetEntry(i, jo)
}

func (s *SortTransformer) Execute(d pimtrace.Data) (pimtrace.Data, error) {
	sort.Sort(&SortTransformerSorter{
		SortTransformer: s,
		Data:            d,
	})
	return d, nil
}

var _ Operation = (*SortTransformer)(nil)

type GroupTransformer struct {
	Columns []*ColumnExpression
}

func (g *GroupTransformer) Execute(d pimtrace.Data) (pimtrace.Data, error) {
	headers := map[string]int{}
	for i, c := range g.Columns {
		headers[c.Name] = i
	}
	td := make([]*groupdata.Row, 0)
	pos := map[string]int{}
	for i := 0; i < d.Len(); i++ {
		r := make([]pimtrace.Value, len(g.Columns))
		e := d.Entry(i)
		for i, c := range g.Columns {
			v, err := c.Operation.Execute(e)
			if err != nil {
				return nil, err
			}
			r[i] = v
		}
		key := pimtrace.SimpleArrayValue(r).String()
		if p, ok := pos[key]; ok {
			td[p].Contents = td[p].Contents.SetEntry(td[p].Contents.Len(), e)
		} else {
			pos[key] = len(td)
			self := d.NewSelf()
			self = self.SetEntry(self.Len(), e)
			td = append(td, &groupdata.Row{
				Headers:  headers,
				Row:      pimtrace.SimpleArrayValue(r),
				Contents: self,
			})
		}
	}
	return groupdata.Data(td), nil
}

var _ Operation = (*GroupTransformer)(nil)
