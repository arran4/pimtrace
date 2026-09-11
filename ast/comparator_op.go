package ast

import (
	"fmt"

	"github.com/arran4/go-evaluator"
)

type SafeComparisonExpression struct {
	Operator string
	Left     evaluator.Term
	Right    evaluator.Term
}

func (c *SafeComparisonExpression) Evaluate(d interface{}, opts ...any) (bool, error) {
	lv, err := c.Left.Evaluate(d, opts...)
	if err != nil {
		return false, err
	}
	rv, err := c.Right.Evaluate(d, opts...)
	if err != nil {
		return false, err
	}

	res, err := evaluator.Compare(lv, rv)
	if err != nil {
		return false, fmt.Errorf("comparison error: %w", err)
	}

	switch c.Operator {
	case "<":
		return res < 0, nil
	case "<=":
		return res <= 0, nil
	case ">":
		return res > 0, nil
	case ">=":
		return res >= 0, nil
	default:
		return false, fmt.Errorf("unknown comparison operator %q", c.Operator)
	}
}
