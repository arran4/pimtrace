package ast

import (
	"fmt"
	"log"
	"pimtrace"
	"pimtrace/dataformats/groupdata"
	"pimtrace/dataformats/tabledata"
	"pimtrace/funcs"
	"sort"
	"strings"
	"unicode"
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

type BooleanExpression interface {
	Execute(d pimtrace.Entry) (bool, error)
}

type NotOp struct {
	Not BooleanExpression
}

func (n *NotOp) Execute(d pimtrace.Entry) (bool, error) {
	v, err := n.Not.Execute(d)
	return !v, err
}

var _ BooleanExpression = (*NotOp)(nil)

type ValueExpression interface {
	funcs.ValueExpression
	ColumnName() string
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
	return d.Get(string(ve)), nil
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
	Op  OpFunc
	LHS ValueExpression
	RHS ValueExpression
}

func (e *Op) Execute(d pimtrace.Entry) (bool, error) {
	if e.LHS == nil {
		return false, fmt.Errorf("LHS invalid issue with Op")
	}
	if e.RHS == nil {
		return false, fmt.Errorf("RHS invalid with Op")
	}
	if e.Op == nil {
		return false, fmt.Errorf("op invalid with Op")
	}
	lhsv, err := e.LHS.Execute(d)
	if err != nil {
		return false, fmt.Errorf("LHS error: %w", err)
	}
	rhsv, err := e.RHS.Execute(d)
	if err != nil {
		return false, fmt.Errorf("RHS error: %w", err)
	}
	return e.Op(rhsv, lhsv)
}

type FilterStatement struct {
	Expression BooleanExpression
}

func (f FilterStatement) Execute(d pimtrace.Data) (pimtrace.Data, error) {
	return Filter(d, f.Expression)
}

var _ Operation = (*FilterStatement)(nil)

type FunctionExpression struct {
	Function string
	Args     []ValueExpression
}

func (fe *FunctionExpression) ColumnName() string {
	elems := []string{fe.Function}
	for _, arg := range fe.Args {
		elems = append(elems, arg.ColumnName())
	}
	return strings.Join(elems, "-")
}

func (fe *FunctionExpression) Execute(d pimtrace.Entry) (pimtrace.Value, error) {
	functions := funcs.Functions[ValueExpression]()
	if f, ok := functions[fe.Function]; ok {
		return f(d, fe.Args)
	}
	return nil, fmt.Errorf("%w: %s", ErrUnknownFunction, fe.Function)
}

var _ ValueExpression = (*FunctionExpression)(nil)

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
		if iv != nil && jv != nil {
			if iv.Equal(jv) {
				continue
			}
			return iv.Less(jv)
		}
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
		log.Printf("Key: %s", key)
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
