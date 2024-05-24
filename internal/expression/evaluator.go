package expression

import "github.com/pkg/errors"

var ErrInvalidExpression = errors.New("invalid expression")

type Evaluator interface {
	Eval(expr string) (float64, error)
}
