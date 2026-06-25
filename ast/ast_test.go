package ast

import (
	"bytes"
	"embed"
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"testing"

	"github.com/arran4/go-evaluator"
	"github.com/google/go-cmp/cmp"
	"fmt"
)

var (
	//go:embed "testdata"
	testdata embed.FS
)

func LoadData1(fn string) pimtrace.Data {
	f, err := testdata.Open(fn)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()
	r, err := tabledata.ReadCSV(f, "test", fn)
	if err != nil {
		panic(err)
	}
	return tabledata.Data(r)
}

func Valueify(ss ...any) (result []pimtrace.Value) {
	for _, s := range ss {
		switch s := s.(type) {
		case string:
			result = append(result, pimtrace.SimpleStringValue(s))
			continue
		case int:
			result = append(result, pimtrace.SimpleIntegerValue(s))
			continue
		default:
			panic("unsupported type")
		}
	}
	return
}

func TestCompoundStatement_Execute(t *testing.T) {
	header1 := map[string]int{"address": 3, "currency": 4, "email": 2, "name": 0, "numberrange": 5, "phone": 1}
	header2 := map[string]int{"Name": 0}
	header3 := map[string]int{"Number": 0, "count": 1, "sum-size": 2}
	header4 := map[string]int{"Count": 2, "Month": 1, "Year": 0}
	tests := []struct {
		name       string
		Statements Operation
		data       pimtrace.Data
		want       pimtrace.Data
		wantErr    bool
	}{
		{
			name: "Simple filter",
			Statements: &FilterStatement{
				Expression: &evaluator.Query{
					Expression: &Op{Op: "eq", LHS: EntryExpression("h.numberrange"), RHS: ConstantExpression("4")},
				},
			},
			data: LoadData1("testdata/data10.csv"),
			want: tabledata.Data{
				{
					Headers: header1,
					Row: Valueify(
						"Jasper Joseph", "(125) 832-4826", "mauris.vestibulum@protonmail.edu",
						"Ap #783-8034 Nunc Street", "$73.44", "4",
					),
				},
			},
			wantErr: false,
		},
		{
			name: "Simple Table column filter",
			Statements: &TableTransformer{
				Columns: []*ColumnExpression{
					{Name: "Name", Operation: EntryExpression("h.name")},
				},
			},
			data: LoadData1("testdata/data10.csv"),
			want: tabledata.Data{
				{Headers: header2, Row: Valueify("Jasper Joseph")},
				{Headers: header2, Row: Valueify("Rogan Hopkins")},
				{Headers: header2, Row: Valueify("Shay Cleveland")},
				{Headers: header2, Row: Valueify("Maite Weaver")},
				{Headers: header2, Row: Valueify("Adria Herring")},
				{Headers: header2, Row: Valueify("Laurel Gonzalez")},
				{Headers: header2, Row: Valueify("Jane Bender")},
				{Headers: header2, Row: Valueify("Melinda Barton")},
				{Headers: header2, Row: Valueify("Colorado Sandoval")},
				{Headers: header2, Row: Valueify("Felix Sutton")},
			},
			wantErr: false,
		},
		{
			name: "Simple Table column filter and sort",
			Statements: &CompoundStatement{Statements: []Operation{
				&TableTransformer{
					Columns: []*ColumnExpression{
						{Name: "Name", Operation: EntryExpression("h.name")},
					},
				},
				&SortTransformer{[]ValueExpression{EntryExpression("c.Name")}},
			}},
			data: LoadData1("testdata/data10.csv"),
			want: tabledata.Data{
				{Headers: header2, Row: Valueify("Adria Herring")},
				{Headers: header2, Row: Valueify("Colorado Sandoval")},
				{Headers: header2, Row: Valueify("Felix Sutton")},
				{Headers: header2, Row: Valueify("Jane Bender")},
				{Headers: header2, Row: Valueify("Jasper Joseph")},
				{Headers: header2, Row: Valueify("Laurel Gonzalez")},
				{Headers: header2, Row: Valueify("Maite Weaver")},
				{Headers: header2, Row: Valueify("Melinda Barton")},
				{Headers: header2, Row: Valueify("Rogan Hopkins")},
				{Headers: header2, Row: Valueify("Shay Cleveland")},
			},
			wantErr: false,
		},
		{
			name: "Summary Table with all the functions and group by number",
			Statements: &CompoundStatement{Statements: []Operation{
				&GroupTransformer{
					Columns: []*ColumnExpression{
						{Name: "numberrange", Operation: EntryExpression("h.numberrange")},
					},
				},
				&TableTransformer{
					Columns: []*ColumnExpression{
						{Name: "Number", Operation: EntryExpression("h.numberrange")},
						{Name: "count", Operation: &FunctionExpression{Function: "count"}}, //Args: []ValueExpression{EntryExpression("t.contents")}}},
						{Name: "sum-size", Operation: &FunctionExpression{Function: "sum", Args: []ValueExpression{EntryExpression("h.numberrange")}}},
					},
				},
			}},
			data: LoadData1("testdata/data10.csv"),
			want: tabledata.Data{
				{
					Headers: header3,
					Row:     Valueify("4", 1, 4*1),
				},
				{
					Headers: header3,
					Row:     Valueify("9", 6, 6*9),
				},
				{
					Headers: header3,
					Row:     Valueify("7", 1, 1*7),
				},
				{
					Headers: header3,
					Row:     Valueify("6", 1, 1*6),
				},
				{
					Headers: header3,
					Row:     Valueify("1", 1, 1*1),
				},
			},
			wantErr: false,
		},
		{
			name: "Summary Table with all the functions and group by formatted columns",
			Statements: &CompoundStatement{Statements: []Operation{
				&GroupTransformer{
					Columns: []*ColumnExpression{
						{Name: "year", Operation: &FunctionExpression{Function: "year", Args: []ValueExpression{EntryExpression("c.date")}}},
						{Name: "month", Operation: &FunctionExpression{Function: "month", Args: []ValueExpression{EntryExpression("c.date")}}},
					},
				},
				&TableTransformer{
					Columns: []*ColumnExpression{
						{Name: "Year", Operation: EntryExpression("c.year")},
						{Name: "Month", Operation: EntryExpression("c.month")},
						{Name: "Count", Operation: &FunctionExpression{Function: "count"}}, //Args: []ValueExpression{EntryExpression("t.contents")}}},
					},
				},
				&SortTransformer{
					Expression: []ValueExpression{
						EntryExpression("c.Year"),
						EntryExpression("c.Month"),
						EntryExpression("c.Count"),
					},
				},
			}},
			data: LoadData1("testdata/tdata1000.csv"),
			want: tabledata.Data{
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2021), pimtrace.SimpleIntegerValue(11), pimtrace.SimpleIntegerValue(74)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2021), pimtrace.SimpleIntegerValue(12), pimtrace.SimpleIntegerValue(70)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(1), pimtrace.SimpleIntegerValue(86)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(2), pimtrace.SimpleIntegerValue(82)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(3), pimtrace.SimpleIntegerValue(89)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(4), pimtrace.SimpleIntegerValue(95)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(5), pimtrace.SimpleIntegerValue(76)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(6), pimtrace.SimpleIntegerValue(81)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(7), pimtrace.SimpleIntegerValue(100)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(8), pimtrace.SimpleIntegerValue(74)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(9), pimtrace.SimpleIntegerValue(78)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(10), pimtrace.SimpleIntegerValue(84)},
				},
				{
					Headers: header4,
					Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(11), pimtrace.SimpleIntegerValue(11)},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.Statements.Execute(tt.data, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Execute() \n%s", diff)
			}
		})
	}
}

func TestSortTransformer_Execute(t *testing.T) {
	yearDateHeader := map[string]int{
		"year-date": 0,
		"count":     1,
	}
	tests := []struct {
		name            string
		SortTransformer SortTransformer
		d               pimtrace.Data
		want            string
		wantErr         bool
	}{
		{
			name:            "Year-Date",
			SortTransformer: SortTransformer{Expression: []ValueExpression{EntryExpression("c.year-date")}},
			d: tabledata.Data{
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2023), pimtrace.SimpleIntegerValue(2190)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2022), pimtrace.SimpleIntegerValue(15664)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2015), pimtrace.SimpleIntegerValue(25332)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2017), pimtrace.SimpleIntegerValue(26803)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2021), pimtrace.SimpleIntegerValue(14606)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{&pimtrace.SimpleNilValue{}, pimtrace.SimpleIntegerValue(175271)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2019), pimtrace.SimpleIntegerValue(31743)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2020), pimtrace.SimpleIntegerValue(26422)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2007), pimtrace.SimpleIntegerValue(24102)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2018), pimtrace.SimpleIntegerValue(24107)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2016), pimtrace.SimpleIntegerValue(29190)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2004), pimtrace.SimpleIntegerValue(9855)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2013), pimtrace.SimpleIntegerValue(33164)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2014), pimtrace.SimpleIntegerValue(38906)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2012), pimtrace.SimpleIntegerValue(33274)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2006), pimtrace.SimpleIntegerValue(17979)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2011), pimtrace.SimpleIntegerValue(39946)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2010), pimtrace.SimpleIntegerValue(34943)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2008), pimtrace.SimpleIntegerValue(14722)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2009), pimtrace.SimpleIntegerValue(28008)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2003), pimtrace.SimpleIntegerValue(9)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2005), pimtrace.SimpleIntegerValue(16396)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(1970), pimtrace.SimpleIntegerValue(2)}},
				{Headers: yearDateHeader, Row: []pimtrace.Value{pimtrace.SimpleIntegerValue(2002), pimtrace.SimpleIntegerValue(2)}},
			},
			want:    "year-date,count\n,175271\n1970,2\n2002,2\n2003,9\n2004,9855\n2005,16396\n2006,17979\n2007,24102\n2008,14722\n2009,28008\n2010,34943\n2011,39946\n2012,33274\n2013,33164\n2014,38906\n2015,25332\n2016,29190\n2017,26803\n2018,24107\n2019,31743\n2020,26422\n2021,14606\n2022,15664\n2023,2190\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.SortTransformer.Execute(tt.d, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			b := bytes.NewBuffer(nil)
			np, ok := got.(pimtrace.CSVOutputCapable)
			if !ok {
				t.Errorf("Not a csv capable data set")
				return
			}
			if err := np.WriteCSVStream(b, "out.csv"); err != nil {
				t.Errorf("Failed CSV write: %s", err)
				return
			}
			if diff := cmp.Diff(b.String(), tt.want); diff != "" {
				t.Errorf("Execute() diff = \n%s", diff)
				t.Errorf("Execute() got = \n%s", b.String())
			}
		})
	}
}

type mockEntry struct {
	vals map[string]pimtrace.Value
}

func (m *mockEntry) Get(key string) (pimtrace.Value, error) {
	if v, ok := m.vals[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found")
}

type filterMockData struct {
	entries []pimtrace.Entry
}

func (m *filterMockData) Len() int { return len(m.entries) }
func (m *filterMockData) Entry(n int) pimtrace.Entry { return m.entries[n] }
func (m *filterMockData) Truncate(n int) pimtrace.Data {
	m.entries = m.entries[:n]
	return m
}
func (m *filterMockData) SetEntry(n int, entry pimtrace.Entry) pimtrace.Data {
	m.entries[n] = entry
	return m
}
func (m *filterMockData) NewSelf() pimtrace.Data { return &filterMockData{} }

func TestFilter(t *testing.T) {
	d := &filterMockData{entries: []pimtrace.Entry{
		&mockEntry{vals: map[string]pimtrace.Value{"keep": pimtrace.SimpleStringValue("yes")}},
		&mockEntry{vals: map[string]pimtrace.Value{"keep": pimtrace.SimpleStringValue("no")}},
		&mockEntry{vals: map[string]pimtrace.Value{"keep": pimtrace.SimpleStringValue("yes")}},
		&mockEntry{vals: map[string]pimtrace.Value{"keep": pimtrace.SimpleStringValue("error")}},
	}}

	// Because of evaluator behavior we just use dummy expressions
	res, err := Filter(d, &evaluator.Query{
		Expression: &dummyBoolExpr{val: true, limit: 2},
	}, nil)

	if err != nil {
		t.Errorf("Filter() error = %v", err)
	}

	if res.Len() != 2 {
		t.Errorf("Filter() expected 2 results, got %d", res.Len())
	}
}

type dummyBoolExpr struct {
	val bool
	limit int
	calls int
}

func (m *dummyBoolExpr) Evaluate(d interface{}, opts ...any) (bool, error) {
	m.calls++
	if m.calls == 4 {
		return false, fmt.Errorf("error case")
	}
	return m.val && m.calls <= m.limit, nil
}

func TestOps(t *testing.T) {
	v1 := pimtrace.SimpleStringValue("hello")
	v2 := pimtrace.SimpleStringValue("hello")
	v3 := pimtrace.SimpleStringValue("world")

	if eq, _ := EqualOp(v1, v2); !eq {
		t.Errorf("EqualOp(v1, v2) = %v, want true", eq)
	}
	if eq, _ := EqualOp(v1, v3); eq {
		t.Errorf("EqualOp(v1, v3) = %v, want false", eq)
	}

	v4 := pimtrace.SimpleStringValue("hello world")
	if contains, _ := ContainsOp(v4, v1); !contains {
		t.Errorf("ContainsOp(v4, v1) = %v, want true", contains)
	}
	if contains, _ := ContainsOp(v1, v4); contains {
		t.Errorf("ContainsOp(v1, v4) = %v, want false", contains)
	}

	v5 := pimtrace.SimpleStringValue("HELLO")
	if icontains, _ := IContainsOp(v1, v5); !icontains {
		t.Errorf("IContainsOp(v1, v5) = %v, want true", icontains)
	}
	if icontains, _ := IContainsOp(v4, v5); !icontains {
		t.Errorf("IContainsOp(v4, v5) = %v, want true", icontains)
	}
}

func TestCompoundStatement_Simplify(t *testing.T) {
	c1 := &CompoundStatement{
		Statements: []Operation{
			&FilterStatement{},
		},
	}
	s1 := c1.Simplify()
	if _, ok := s1.(*FilterStatement); !ok {
		t.Errorf("Simplify() of 1 len CompoundStatement failed")
	}

	c0 := &CompoundStatement{
		Statements: []Operation{},
	}
	if s0 := c0.Simplify(); s0 != nil {
		t.Errorf("Simplify() of 0 len CompoundStatement failed")
	}

	c2 := &CompoundStatement{
		Statements: []Operation{
			&FilterStatement{},
			&CompoundStatement{
				Statements: []Operation{
					&FilterStatement{},
				},
			},
		},
	}
	s2 := c2.Simplify()
	if cs, ok := s2.(*CompoundStatement); !ok || len(cs.Statements) != 2 {
		t.Errorf("Simplify() of nested CompoundStatement failed")
	}
}

func TestValueExpressions(t *testing.T) {
	c := ConstantExpression("test_value")
	if n := c.ColumnName(); n != "test-value" {
		t.Errorf("ConstantExpression.ColumnName() = %v, want test-value", n)
	}

	e := EntryExpression("c.name")
	if n := e.ColumnName(); n != "name" {
		t.Errorf("EntryExpression.ColumnName() = %v, want name", n)
	}

	f := &FunctionExpression{
		Function: "count",
		Args: []ValueExpression{
			e,
		},
	}
	if n := f.ColumnName(); n != "count-name" {
		t.Errorf("FunctionExpression.ColumnName() = %v, want count-name", n)
	}

	ef := &EvaluatorFunctionExpression{
		Function: "sum",
		FunctionExpression: evaluator.FunctionExpression{
			Name: "sum",
			Args: []evaluator.Term{
				e,
			},
		},
	}
	if n := ef.ColumnName(); n != "sum-name" {
		t.Errorf("EvaluatorFunctionExpression.ColumnName() = %v, want sum-name", n)
	}
}

func TestToPimtraceValue(t *testing.T) {
	if v, _ := toPimtraceValue(nil); v.Type() != pimtrace.Nil {
		t.Errorf("toPimtraceValue(nil) failed")
	}

	if v, _ := toPimtraceValue(pimtrace.SimpleStringValue("test")); v.Type() != pimtrace.String {
		t.Errorf("toPimtraceValue(pimtrace.Value) failed")
	}

	if v, _ := toPimtraceValue(int(42)); v.Type() != pimtrace.Integer {
		t.Errorf("toPimtraceValue(int) failed")
	}

	if v, _ := toPimtraceValue(int64(42)); v.Type() != pimtrace.Integer {
		t.Errorf("toPimtraceValue(int64) failed")
	}

	if v, _ := toPimtraceValue(float64(42.5)); v.Type() != pimtrace.Integer {
		t.Errorf("toPimtraceValue(float64) failed")
	}

	if v, _ := toPimtraceValue("hello"); v.Type() != pimtrace.String {
		t.Errorf("toPimtraceValue(string) failed")
	}

	if v, _ := toPimtraceValue([]interface{}{"test", 42}); v.Type() != pimtrace.Array {
		t.Errorf("toPimtraceValue([]interface{}) failed")
	}

	if _, err := toPimtraceValue([]interface{}{struct{}{}}); err == nil {
		t.Errorf("toPimtraceValue([]interface{} with bad inner) should error")
	}

	if _, err := toPimtraceValue(struct{}{}); err == nil {
		t.Errorf("toPimtraceValue(unsupported) should error")
	}
}

func TestEntryExpression_Execute(t *testing.T) {
	ve := EntryExpression("c.val")
	ctx := &evaluator.Context{}
	d := &mockEntry{vals: map[string]pimtrace.Value{"c.val": pimtrace.SimpleIntegerValue(42)}}

	res, err := ve.Execute(d, ctx)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
	if iv, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(iv) != 42 {
		t.Errorf("Execute() result = %v, want 42", res)
	}

	evalRes, err := ve.Evaluate(evaluatorEntryWrapper{Entry: d})
	if err != nil {
		t.Errorf("Evaluate() error = %v", err)
	}
	if iv, ok := evalRes.(pimtrace.SimpleIntegerValue); !ok || int(iv) != 42 {
		t.Errorf("Evaluate() result = %v, want 42", evalRes)
	}
}

func TestConstantExpression_Execute(t *testing.T) {
	ve := ConstantExpression("test")
	ctx := &evaluator.Context{}
	res, err := ve.Execute(nil, ctx)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
	if sv, ok := res.(pimtrace.SimpleStringValue); !ok || string(sv) != "test" {
		t.Errorf("Execute() result = %v, want test", res)
	}

	evalRes, err := ve.Evaluate(nil)
	if err != nil {
		t.Errorf("Evaluate() error = %v", err)
	}
	if sv, ok := evalRes.(pimtrace.SimpleStringValue); !ok || string(sv) != "test" {
		t.Errorf("Evaluate() result = %v, want test", evalRes)
	}
}

type mockErrOp struct{}

func (m *mockErrOp) Execute(d pimtrace.Data, ctx *evaluator.Context) (pimtrace.Data, error) {
	return nil, fmt.Errorf("mock error")
}

func TestCompoundStatement_ExecuteError(t *testing.T) {
	c := &CompoundStatement{
		Statements: []Operation{
			&mockErrOp{},
		},
	}
	_, err := c.Execute(nil, nil)
	if err == nil {
		t.Errorf("CompoundStatement.Execute error expected")
	}
}

func TestFunctionExpression_ExecuteError(t *testing.T) {
	fe := &FunctionExpression{Function: "unknown_func"}
	_, err := fe.Execute(nil, nil)
	if err == nil {
		t.Errorf("Execute() unknown func expected error")
	}
}

// A mock struct for evaluator.Function
type mockEvalFunc struct {
	callFunc func(args ...interface{}) (interface{}, error)
}

func (m *mockEvalFunc) Call(args ...interface{}) (interface{}, error) {
	return m.callFunc(args...)
}

func TestEvaluatorFunctionExpression_Execute(t *testing.T) {
	ef := &EvaluatorFunctionExpression{
		Function: "sum",
		FunctionExpression: evaluator.FunctionExpression{
			Name: "sum",
			Args: []evaluator.Term{
				EntryExpression("c.val"),
			},
		},
	}

	ctx := &evaluator.Context{
		Functions: map[string]evaluator.Function{
			"sum": &mockEvalFunc{
				callFunc: func(args ...interface{}) (interface{}, error) {
					return 42, nil
				},
			},
		},
	}

	d := &mockEntry{vals: map[string]pimtrace.Value{"c.val": pimtrace.SimpleIntegerValue(10)}}

	res, err := ef.Execute(d, ctx)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if iv, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(iv) != 42 {
		t.Errorf("Execute() result = %v, want 42", res)
	}

	// Test non-Value return from function that toPimtraceValue can convert
	ctx.Functions["sum"] = &mockEvalFunc{
		callFunc: func(args ...interface{}) (interface{}, error) {
			return "hello", nil
		},
	}
	res2, err := ef.Execute(d, ctx)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
	if sv, ok := res2.(pimtrace.SimpleStringValue); !ok || string(sv) != "hello" {
		t.Errorf("Execute() result = %v, want hello", res2)
	}

	// Test error return
	ctx.Functions["sum"] = &mockEvalFunc{
		callFunc: func(args ...interface{}) (interface{}, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	_, err = ef.Execute(d, ctx)
	if err == nil {
		t.Errorf("Execute() expected error, got nil")
	}
}

func TestFunctionExpression_Evaluate(t *testing.T) {
	fe := &FunctionExpression{Function: "unknown_func"}

	d := &mockEntry{}
	w := evaluatorEntryWrapper{Entry: d}

	ctx := &evaluator.Context{}
	_, err := fe.Evaluate(w, ctx)
	if err == nil {
		t.Errorf("Evaluate() error expected for unknown func")
	}

	// invalid entry type
	_, err = fe.Evaluate(nil)
	if err == nil {
		t.Errorf("Evaluate() invalid type expected error")
	}
}

type valueEntry struct {
	pimtrace.SimpleStringValue
}

func (v valueEntry) Get(k string) (pimtrace.Value, error) {
	return v.SimpleStringValue, nil
}

func TestEntryPathor_Find(t *testing.T) {
	d := &mockEntry{vals: map[string]pimtrace.Value{
		"c.val": pimtrace.SimpleIntegerValue(42),
		"c.sub": valueEntry{pimtrace.SimpleStringValue("subval")},
	}}
	ep := NewEntryPathor(d)

	// Test empty path
	if p := ep.Find(""); p != ep {
		t.Errorf("Find empty path should return self")
	}

	// Test successful Get
	p := ep.Find("c.val")
	if p == nil {
		t.Errorf("Find(c.val) should not be nil")
	}

	// Test successful Get returning Entry
	p3 := ep.Find("c.sub")
	if _, ok := p3.(*EntryPathor); !ok {
		t.Errorf("Find(c.sub) returning entry should wrap in EntryPathor")
	}

	// Test unsuccessful Get
	p2 := ep.Find("nonexistent")
	if p2 == nil {
		t.Errorf("Find(nonexistent) should not be nil but Invalidor")
	}
}
