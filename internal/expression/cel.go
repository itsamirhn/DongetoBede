package expression

import (
	"github.com/google/cel-go/cel"
)

type celEvaluator struct {
	env *cel.Env
}

func NewCelEvaluator() Evaluator {
	env, err := cel.NewEnv()
	if err != nil {
		panic(err)
	}
	return &celEvaluator{
		env: env,
	}
}

func (c celEvaluator) Eval(expr string) (float64, error) {
	ast, iss := c.env.Parse(expr)
	if iss != nil && iss.Err() != nil {
		return 0, iss.Err()
	}
	prg, err := c.env.Program(ast)
	if err != nil {
		return 0, err
	}
	out, _, err := prg.Eval(map[string]interface{}{})
	if err != nil {
		return 0, err
	}
	switch out.Type() {
	case cel.DoubleType:
		return out.Value().(float64), nil
	case cel.UintType:
		return float64(out.Value().(uint)), nil
	case cel.IntType:
		return float64(out.Value().(int64)), nil
	default:
		return 0, ErrInvalidExpression
	}
}
