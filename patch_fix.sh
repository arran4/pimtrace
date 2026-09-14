cat << 'INNER_EOF' > /tmp/fix.diff
--- a/ast/ast_test.go
+++ b/ast/ast_test.go
@@ -972,8 +972,11 @@ func TestSortTransformer_Execute_AggregateCountDescending(t *testing.T) {
	stDesc := SortTransformer{Keys: []SortKey{{Expression: EntryExpression("c.count"), Direction: Descending}}}
	resDesc, _ := stDesc.Execute(d, nil)

	if *resDesc.Entry(0).(*tabledata.Row).Row[0].Integer() != 20 {
+		t.Errorf("Expected 20 first in Descending count sort, got %v", resDesc.Entry(0).(*tabledata.Row).Row[0])
	}
	if *resDesc.Entry(2).(*tabledata.Row).Row[0].Integer() != 5 {
+		t.Errorf("Expected 5 last in Descending count sort, got %v", resDesc.Entry(2).(*tabledata.Row).Row[0])
	}
 }
INNER_EOF
patch -p1 < /tmp/fix.diff
