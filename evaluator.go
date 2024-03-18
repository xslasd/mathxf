package mathxf

import (
	"fmt"
	"math"
)

type IEvaluator interface {
	INode
	GetPositionToken() *Token
	Evaluate(ctx EvaluatorContext) (*Value, error)
}
type functionCallArgument interface {
	Evaluate(ctx EvaluatorContext) (*Value, error)
}

// Expression 处理TokenAnd 和 TokenOr
type Expression struct {
	expr1   IEvaluator
	expr2   IEvaluator
	opToken *Token
}

func (e Expression) Execute(ctx EvaluatorContext) error {
	_, err := e.expr1.Evaluate(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e Expression) GetPositionToken() *Token {
	return e.expr1.GetPositionToken()
}

func (e Expression) Evaluate(ctx EvaluatorContext) (*Value, error) {
	v1, err := e.expr1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if e.expr2 != nil {
		switch e.opToken.typ {
		case TokenAnd:
			if !v1.IsTrue() {
				return AsValue(false), nil
			} else {
				v2, err := e.expr2.Evaluate(ctx)
				if err != nil {
					return nil, err
				}
				return AsValue(v2.IsTrue()), nil
			}
		case TokenOr:
			if v1.IsTrue() {
				return AsValue(true), nil
			} else {
				v2, err := e.expr2.Evaluate(ctx)
				if err != nil {
					return nil, err
				}
				return AsValue(v2.IsTrue()), nil
			}
		default:
			return nil, nil //ctx.Error(fmt.Sprintf("unimplemented: %name", r.opToken.Val), r.opToken)
		}
	} else {
		return v1, nil
	}
}

// relationalExpression 处理  TokenEqual  TokenNotEqual  TokenLess  TokenLessEqual  TokenGreater  TokenGreaterEqual
type relationalExpression struct {
	expr1   IEvaluator
	expr2   IEvaluator
	opToken *Token
}

func (r relationalExpression) Execute(ctx EvaluatorContext) error {
	v, err := r.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(v, "---------------", r.opToken.val)
	return nil
}

func (r relationalExpression) GetPositionToken() *Token {
	return r.expr1.GetPositionToken()
}

func (r relationalExpression) Evaluate(ctx EvaluatorContext) (*Value, error) {
	v1, err := r.expr1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if r.expr2 != nil {
		v2, err := r.expr2.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		switch r.opToken.typ {
		case TokenLessEquals:
			if v1.IsFloat() || v2.IsFloat() {
				return AsValue(v1.Float() <= v2.Float()), nil
			}
			if v1.IsTime() && v2.IsTime() {
				tm1, tm2 := v1.Time(), v2.Time()
				return AsValue(tm1.Before(tm2) || tm1.Equal(tm2)), nil
			}
			return AsValue(v1.Integer() <= v2.Integer()), nil
		case TokenGreatEquals:
			if v1.IsFloat() || v2.IsFloat() {
				return AsValue(v1.Float() >= v2.Float()), nil
			}
			if v1.IsTime() && v2.IsTime() {
				tm1, tm2 := v1.Time(), v2.Time()
				return AsValue(tm1.After(tm2) || tm1.Equal(tm2)), nil
			}
			return AsValue(v1.Integer() >= v2.Integer()), nil
		case TokenEquals:
			return AsValue(v1.EqualValueTo(v2)), nil
		case TokenGreat:
			if v1.IsFloat() || v2.IsFloat() {
				return AsValue(v1.Float() > v2.Float()), nil
			}
			if v1.IsTime() && v2.IsTime() {
				return AsValue(v1.Time().After(v2.Time())), nil
			}
			return AsValue(v1.Integer() > v2.Integer()), nil
		case TokenLess:
			if v1.IsFloat() || v2.IsFloat() {
				return AsValue(v1.Float() < v2.Float()), nil
			}
			if v1.IsTime() && v2.IsTime() {
				return AsValue(v1.Time().Before(v2.Time())), nil
			}
			return AsValue(v1.Integer() < v2.Integer()), nil
		case TokenNotEquals:
			return AsValue(!v1.EqualValueTo(v2)), nil
		case TokenIn:
			return AsValue(v2.Contains(v1)), nil
		default:
			return nil, fmt.Errorf("unimplemented: %v", r.opToken.val)
		}
	} else {
		return v1, nil
	}
}

// simpleExpression 处理 TokenAdd TokenSub
type simpleExpression struct {
	term1   IEvaluator
	term2   IEvaluator
	opToken *Token
}

func (s simpleExpression) Execute(ctx EvaluatorContext) error {
	v, err := s.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%.16f\n", v.Float())
	return nil
}

func (s simpleExpression) GetPositionToken() *Token {
	return s.term1.GetPositionToken()
}

func (s simpleExpression) Evaluate(ctx EvaluatorContext) (*Value, error) {
	t1, err := s.term1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	result := t1
	if s.term2 != nil {
		t2, err := s.term2.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		switch s.opToken.typ {
		case TokenAdd:
			if result.IsString() || t2.IsString() {
				// Result will be a string
				return AsValue(result.String() + t2.String()), nil
			}
			if result.IsFloat() || t2.IsFloat() {
				// Result will be a float
				return AsValue(result.Float() + t2.Float()), nil
			}
			// Result will be an integer
			return AsValue(result.Integer() + t2.Integer()), nil
		case TokenSub:
			if result.IsFloat() || t2.IsFloat() {
				// Result will be a float
				return AsValue(result.Float() - t2.Float()), nil
			}
			// Result will be an integer
			return AsValue(result.Integer() - t2.Integer()), nil
		default:
			return nil, nil //ctx.Error("Unimplemented", name.GetPositionToken())
		}
	}

	return result, nil
}

// termExpression 处理 TokenMul TokenDiv TokenMod
type termExpression struct {
	factor1 IEvaluator
	factor2 IEvaluator
	opToken *Token
}

func (t termExpression) Execute(ctx EvaluatorContext) error {
	v, err := t.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}

func (t termExpression) GetPositionToken() *Token {
	return t.factor1.GetPositionToken()
}

func (t termExpression) Evaluate(ctx EvaluatorContext) (*Value, error) {
	f1, err := t.factor1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if t.factor2 != nil {
		f2, err := t.factor2.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		switch t.opToken.typ {
		case TokenMul:
			if f1.IsFloat() || f2.IsFloat() {
				// Result will be float
				return AsValue(f1.Float() * f2.Float()), nil
			}
			// Result will be int
			return AsValue(f1.Integer() * f2.Integer()), nil
		case TokenDiv:
			if f1.IsFloat() || f2.IsFloat() {
				// Result will be float
				divisor := f2.Float()
				if divisor == 0 {
					return nil, nil //ctx.Error("float divide by zero", t.factor2.GetPositionToken())
				}
				return AsValue(f1.Float() / divisor), nil
			}
			// Result will be int
			divisor := f2.Integer()
			if divisor == 0 {
				return nil, nil // ctx.Error("integer divide by zero", t.factor2.GetPositionToken())
			}
			return AsValue(f1.Integer() / divisor), nil
		case TokenMod:
			// Result will be int
			divisor := f2.Integer()
			if divisor == 0 {
				return nil, nil // ctx.Error("integer divide by zero", t.factor2.GetPositionToken())
			}
			return AsValue(f1.Integer() % divisor), nil
		default:
			return nil, nil // ctx.Error("unimplemented", t.opToken)
		}
	} else {
		return f1, nil
	}
}

// powerExpression 处理 returns x**y, the base-x exponential of y.
type powerExpression struct {
	power1 IEvaluator
	power2 IEvaluator
}

func (p powerExpression) Execute(ctx EvaluatorContext) error {
	v, err := p.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}

func (p powerExpression) GetPositionToken() *Token {
	return p.power1.GetPositionToken()
}

func (p powerExpression) Evaluate(ctx EvaluatorContext) (*Value, error) {
	p1, err := p.power1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if p.power2 != nil {
		p2, err := p.power2.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		return AsValue(math.Pow(p1.Float(), p2.Float())), nil
	}
	return p1, nil
}
