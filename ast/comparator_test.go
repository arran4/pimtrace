package ast

import (
	"pimtrace"
	"testing"
	"time"

)

func TestComparatorAdapter(t *testing.T) {
	t.Run("CSV numeric range integration equivalent", func(t *testing.T) {
		ca := &ComparatorAdapter{Value: pimtrace.SimpleStringValue("2")}

		// 2 < 10
		res, err := ca.Compare(10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != -1 {
			t.Errorf("expected -1, got %d", res)
		}

		// 2 < 10.25 (float)
		res, err = ca.Compare(10.25)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != -1 {
			t.Errorf("expected -1, got %d", res)
		}

		// negative values
		caNeg := &ComparatorAdapter{Value: pimtrace.SimpleStringValue("-5")}
		res, err = caNeg.Compare(-2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != -1 {
			t.Errorf("expected -1 for -5 < -2, got %d", res)
		}

		// zero
		caZero := &ComparatorAdapter{Value: pimtrace.SimpleStringValue("0")}
		res, err = caZero.Compare(0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != 0 {
			t.Errorf("expected 0 for 0 == 0, got %d", res)
		}
	})

	t.Run("Mail text boolean combination", func(t *testing.T) {
		ca := &ComparatorAdapter{Value: pimtrace.SimpleStringValue("test")}
		res, err := ca.Compare("test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != 0 {
			t.Errorf("expected 0, got %d", res)
		}
	})

	t.Run("iCal date range example", func(t *testing.T) {
		ca := &ComparatorAdapter{Value: pimtrace.SimpleStringValue("2023-01-01T12:00:00Z")}
		otherTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		res, err := ca.Compare(otherTime)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != -1 {
			t.Errorf("expected -1, got %d", res)
		}
	})

	t.Run("Invalid coercion error", func(t *testing.T) {
		ca := &ComparatorAdapter{Value: pimtrace.SimpleStringValue("not a number")}
		_, err := ca.Compare(10)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Pimtrace value vs Pimtrace value", func(t *testing.T) {
		ca1 := &ComparatorAdapter{Value: pimtrace.SimpleStringValue("10")}
		ca2 := &ComparatorAdapter{Value: pimtrace.SimpleStringValue("2")}
		// String comparison 10 < 2 is true lexicographically
		res, err := ca1.Compare(ca2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != -1 {
			t.Errorf("expected -1 for \"10\" < \"2\", got %d", res)
		}
	})
}
